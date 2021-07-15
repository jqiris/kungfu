package jobs

import (
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/utils"
	"sort"
	"sync"
	"time"
)

var (
	keeper *JobKeeper
)

func InitJobs() {
	keeper = NewJobKeeper()
	keeper.ExecJob()
}

func AddJob(delay time.Duration, job JobWorker) {
	go keeper.AddJob(delay, job)
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
	AddTime   time.Time //添加时间
	StartTime time.Time //开始时间
	Worker    JobWorker //任务对象
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
	finishTime := time.Now()
	logger.Infof(
		"job finished,name:%v,addtime:%v,starttime:%v,endtime:%v, total:%v秒, deal:%v秒",
		worker.Name(),
		s.AddTime.Format("2006-01-02 15:04:05"),
		s.StartTime.Format("2006-01-02 15:04:05"),
		finishTime.Format("2006-01-02 15:04:05"),
		finishTime.Sub(s.AddTime).Seconds(),
		finishTime.Sub(s.StartTime).Seconds(),
	)
}

func NewJobItem(sTime time.Time, worker JobWorker) *JobItem {
	return &JobItem{
		AddTime:   time.Now(),
		StartTime: sTime,
		Worker:    worker,
	}
}

type JobQueue struct {
	StartTime time.Time     //开始时间
	JobItems  *Queue        //任务队列
	mutex     *sync.RWMutex //锁
}

func NewJobQueue(sTime time.Time) *JobQueue {
	return &JobQueue{
		StartTime: sTime,
		JobItems:  NewQueue(),
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
	s.JobItems.RangePop(func(item interface{}) bool {
		if job, ok := item.(*JobItem); ok {
			if job != nil {
				go utils.SafeRun(func() {
					if job.Worker != nil {
						job.Worker.BeforeExec()
					}
					job.ExecJob()
				})
			}
		}
		return true
	})
}

type JobQueues []*JobQueue

func (q JobQueues) Len() int           { return len(q) }
func (q JobQueues) Less(i, j int) bool { return q[i].StartTime.Unix() < q[j].StartTime.Unix() }
func (q JobQueues) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

type JobKeeper struct {
	List  JobQueues
	Index map[time.Time]*JobQueue
	mutex *sync.RWMutex
}

func NewJobKeeper() *JobKeeper {
	return &JobKeeper{
		List:  make(JobQueues, 0),
		Index: make(map[time.Time]*JobQueue),
		mutex: new(sync.RWMutex),
	}
}

func (k *JobKeeper) AddJob(delay time.Duration, job JobWorker) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	sTime := time.Now().Add(delay)
	jobItem := NewJobItem(sTime, job)
	if q, ok := k.Index[sTime]; ok {
		q.AddJob(jobItem)
	} else {
		qs := NewJobQueue(sTime)
		qs.AddJob(jobItem)
		k.Index[sTime] = qs
		k.List = append(k.List, qs)
		sort.Sort(k.List)
	}
}

func (k *JobKeeper) ExecJob() {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				k.mutex.RLock()
				if k.List.Len() > 0 && k.List[0].StartTime.Before(time.Now()) {
					tmp := k.List[0]
					k.List = k.List[1:]
					delete(k.Index, tmp.StartTime)
					tmp.ExeJob()
				}
				k.mutex.RUnlock()
			}
		}
	}()
}
