package jobs

import (
	"sort"
	"sync"
	"time"

	"github.com/jqiris/kungfu/v2/ds"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/utils"
)

var (
	keeper *JobKeeper
)

func init() {
	keeper = NewJobKeeper()
	keeper.ExecJob()
}

func AddJob(delay time.Duration, job JobWorker, options ...ItemOption) {
	go keeper.AddJob(delay, job, options...)
}

func DelJob(id int64) {
	keeper.DelJob(id)
}

type JobWorker interface {
	Name() string                    //任务名称
	BeforeExec()                     //任务执行前操作
	CanExec() bool                   //是否可以执行
	JobExec() bool                   //执行任务,返回执行是否完成
	FailNext() (bool, time.Duration) //失败是否继续执行并且延缓执行时间
	JobFinish()                      //任务执行完成触发
}

type JobItem struct {
	JobId     int64     //任务标识
	AddTime   int64     //添加时间
	StartTime int64     //开始时间
	Worker    JobWorker //任务对象
	Debug     bool      //是否调试
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
				time.AfterFunc(delay, s.ExecJob)
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
		finishTime := time.Now().Unix()
		logger.Infof(
			"job finished,name:%v,addtime:%v,starttime:%v,endtime:%v, total:%v秒, deal:%v秒",
			worker.Name(),
			s.AddTime,
			s.StartTime,
			finishTime,
			finishTime-s.AddTime,
			finishTime-s.StartTime,
		)
	}
}

func NewJobItem(delay time.Duration, worker JobWorker, options ...ItemOption) *JobItem {
	nowTime := time.Now()
	startTime := nowTime.Add(delay)
	job := &JobItem{
		AddTime:   nowTime.Unix(),
		StartTime: startTime.Unix(),
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
			return job.JobId == delId
		}
		return false
	})
}

type JobQueues []*JobQueue

func (q JobQueues) Len() int           { return len(q) }
func (q JobQueues) Less(i, j int) bool { return q[i].StartTime < q[j].StartTime }
func (q JobQueues) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

type JobKeeper struct {
	List    JobQueues
	Index   map[int64]*JobQueue
	AddChan chan *JobItem
	IdList  map[int64]map[int64]int
	DelChan chan int64 //删除id任务
}

func NewJobKeeper() *JobKeeper {
	return &JobKeeper{
		List:    make(JobQueues, 0),
		Index:   make(map[int64]*JobQueue),
		IdList:  make(map[int64]map[int64]int),
		AddChan: make(chan *JobItem, 20),
		DelChan: make(chan int64, 20),
	}
}

func (k *JobKeeper) AddJob(delay time.Duration, job JobWorker, options ...ItemOption) {
	jobItem := NewJobItem(delay, job, options...)
	k.AddChan <- jobItem
}
func (k *JobKeeper) DelJob(id int64) {
	k.DelChan <- id
}

func (k *JobKeeper) ExecJob() {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				nowUnix := time.Now().Unix()
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
			case delId := <-k.DelChan:
				if list, ok := k.IdList[delId]; ok {
					for sTime := range list {
						if q, ok := k.Index[sTime]; ok {
							q.DelJob(delId)
						}
					}
				}
			}
		}
	}()
}
