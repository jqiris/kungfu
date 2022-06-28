package jobs

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jqiris/kungfu/v2/ds"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/utils"
)

var (
	keeper *JobKeeper
)

func init() {
	keeper = NewJobKeeper("default")
	keeper.ExecJob()
}

func AddJob(delay time.Duration, job JobWorker, options ...ItemOption) {
	go keeper.AddJob(delay, job, options...)
}

func DelJob(id int64) {
	keeper.DelJob(id)
}

type JobWorker interface {
	String() string                  //任务名称
	BeforeExec()                     //任务执行前操作
	CanExec() bool                   //是否可以执行
	JobExec() bool                   //执行任务,返回执行是否完成
	FailNext() (bool, time.Duration) //失败是否继续执行并且延缓执行时间
	JobFinish()                      //任务执行完成触发
}

type JobItem struct {
	JobId     int64     //任务标识
	AddTime   int64     //添加时间
	BeginTime int64     //开始时间
	Worker    JobWorker //任务对象
	StartTime int64     //开始时间
	Debug     bool      //是否调试
	Id        string    //任务唯一标识
}

func (s *JobItem) ExecJob() {
	worker := s.Worker
	if worker == nil {
		return
	}
	if worker.CanExec() {
		if worker.JobExec() {
			s.FinishJob()
		} else {
			if next, delay := worker.FailNext(); next {
				s.StartTime = time.Now().Add(delay).UnixMilli()
				keeper.AddChan <- s
			} else {
				s.FinishJob()
			}
		}
	} else {
		s.FinishJob()
	}
}
func (s *JobItem) FinishJob() {
	worker := s.Worker
	if worker == nil {
		return
	}
	worker.JobFinish()
	if s.Debug {
		finishTime := time.Now().UnixMilli()
		logger.Infof(
			"job finished,JobId:%v, AddTime:%v, BeginTime:%v, Worker:%v, EndTime:%v, Total:%v毫秒, Deal:%v毫秒, Id:%v",
			s.JobId,
			s.AddTime,
			s.BeginTime,
			worker,
			finishTime,
			finishTime-s.AddTime,
			finishTime-s.BeginTime,
			s.Id,
		)
	}
}

func NewJobItem(delay time.Duration, worker JobWorker, options ...ItemOption) *JobItem {
	nowTime := time.Now()
	startTime := nowTime.Add(delay)
	job := &JobItem{
		Id:        uuid.NewString(),
		AddTime:   nowTime.UnixMilli(),
		StartTime: startTime.UnixMilli(),
		BeginTime: startTime.UnixMilli(),
		Worker:    worker,
	}
	for _, option := range options {
		option(job)
	}
	return job
}

type JobQueue struct {
	StartTime int64         //开始时间
	JobItems  *ds.Queue     //任务队列
	mutex     *sync.RWMutex //锁
}

func NewJobQueue(sTime int64) *JobQueue {
	return &JobQueue{
		StartTime: sTime,
		JobItems:  ds.NewQueue(),
		mutex:     new(sync.RWMutex),
	}
}

func (s *JobQueue) AddJob(job *JobItem) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.JobItems.Push(job)
	if job.Debug {
		logger.Infof("add job %+v", job)
	}
}

func (s *JobQueue) ExeJob() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.JobItems.RangePop(func(item any) bool {
		if job, ok := item.(*JobItem); ok && job != nil {
			go utils.SafeRun(func() {
				if job.Worker != nil {
					job.Worker.BeforeExec()
				}
				job.ExecJob()
			})
		}
		return true
	})
}
func (s *JobQueue) DelJob(delId int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.JobItems.RangePop(func(item any) bool {
		if job, ok := item.(*JobItem); ok && job != nil {
			canDelete := job.JobId == delId
			if canDelete && job.Debug {
				logger.Infof("delete job %+v", job)
			}
			return canDelete
		}
		return false
	})
}

type JobQueues []*JobQueue

func (q JobQueues) Len() int           { return len(q) }
func (q JobQueues) Less(i, j int) bool { return q[i].StartTime < q[j].StartTime }
func (q JobQueues) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

type JobKeeper struct {
	Name     string
	List     JobQueues
	Index    map[int64]*JobQueue
	AddChan  chan *JobItem
	DelChan  chan int64
	StopChan chan struct{}
	IdList   map[int64]map[int64]int
}

func NewJobKeeper(name string) *JobKeeper {
	return &JobKeeper{
		Name:     name,
		List:     make(JobQueues, 0),
		Index:    make(map[int64]*JobQueue),
		IdList:   make(map[int64]map[int64]int),
		AddChan:  make(chan *JobItem, 20),
		DelChan:  make(chan int64, 20),
		StopChan: make(chan struct{}, 1),
	}
}

func (k *JobKeeper) AddJob(delay time.Duration, job JobWorker, options ...ItemOption) {
	jobItem := NewJobItem(delay, job, options...)
	k.AddChan <- jobItem
}

func (k *JobKeeper) DelJob(delId int64) {
	k.DelChan <- delId
}

func (k *JobKeeper) delJob(delId int64) {
	if list, ok := k.IdList[delId]; ok {
		for sTime := range list {
			if q, ok := k.Index[sTime]; ok {
				q.DelJob(delId)
			}
		}
	}
	delete(k.IdList, delId)
}

func (k *JobKeeper) Stop() {
	k.StopChan <- struct{}{}
}

func (k *JobKeeper) ExecJob() {
	go utils.SafeRun(func() {
		for {
			select {
			case delId := <-k.DelChan:
				k.delJob(delId)
			case now := <-time.After(100 * time.Millisecond):
				nowUnix := now.UnixMilli()
				var jobQueue *JobQueue
				index, left := -1, false
				for index, jobQueue = range k.List {
					if jobQueue.StartTime > nowUnix {
						left = true
						break
					}
					go jobQueue.ExeJob()
					delete(k.Index, jobQueue.StartTime)
				}
				if left {
					k.List = k.List[index:]
				} else {
					k.List = []*JobQueue{}
				}
			case jobItem := <-k.AddChan:
			priority:
				for {
					select {
					case delId := <-k.DelChan:
						k.delJob(delId)
					default:
						break priority
					}
				}
				sTime := jobItem.StartTime
				if q, ok := k.Index[sTime]; ok {
					q.AddJob(jobItem)
				} else {
					qs := NewJobQueue(sTime)
					qs.AddJob(jobItem)
					k.Index[sTime] = qs
					k.List = append(k.List, qs)
					sort.Sort(k.List)
				}
				if jobItem.JobId > 0 {
					if _, ok := k.IdList[jobItem.JobId]; !ok {
						k.IdList[jobItem.JobId] = make(map[int64]int)
					}
					k.IdList[jobItem.JobId][sTime]++
				}
			case <-k.StopChan:
				logger.Infof("job keeper %v received stop signal", k.Name)
				return
			}
		}
	})
}
