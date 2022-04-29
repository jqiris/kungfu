package jobs

import (
	"time"
)

type Job struct {
	name     string        //任务名称
	interval time.Duration //任务间隔
	repeat   int           //重复次数
	count    int           //尝试次数
	worker   func()        //执行函数
}

func NewJob(worker func(), options ...Option) *Job {
	job := &Job{
		worker: worker,
		count:  0,
		repeat: 1,
		name:   "job",
	}
	for _, option := range options {
		option(job)
	}
	if job.repeat > 1 && job.interval == 0 {
		job.interval = time.Second
	}
	return job
}

func (j *Job) String() string {
	return j.name
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) BeforeExec() {
}

func (j *Job) CanExec() bool {
	return j.count < j.repeat
}

func (j *Job) JobExec() bool {
	j.count++
	j.worker()
	return !j.CanExec()
}

func (j *Job) FailNext() (bool, time.Duration) {
	return j.CanExec(), j.interval
}

func (j *Job) JobFinish() {
}
