package jobs

import (
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"testing"
	"time"
)

type JobTest struct {
	name  string
	times int
}

func (j *JobTest) Name() string {
	return "测试任务:" + j.name
}

func (j *JobTest) BeforeExec() {
	fmt.Println("测试任务:" + j.name + ",准备执行")
}

func (j *JobTest) CanExec() bool {
	if j.times < 5 {
		return true
	}
	return false
}

func (j *JobTest) JobExec() bool {
	j.times++
	fmt.Printf("时间:%v,第%v次执行任务:%v \n", time.Now(), j.times, j.name)
	if j.times >= 3 {
		return true
	}
	return false
}

func (j *JobTest) FailNext() (bool, time.Duration) {
	return j.CanExec(), 1 * time.Second
}

func (j *JobTest) JobFinish() {
	fmt.Printf("时间:%v,完成任务：%v \n", time.Now(), j.name)
}

func TestJobs(t *testing.T) {
	logger.Infof("now begin:")
	keeper := NewJobKeeper()
	go func() {
		for i := 1; i < 20; i++ {
			jobn := &JobTest{name: fmt.Sprintf("子任务：%v", i)}
			keeper.AddJob(time.Duration(i)*time.Second, jobn)
		}
	}()
	job3 := &JobTest{name: "任务3"}
	keeper.AddJob(3*time.Second, job3)

	job4 := &JobTest{name: "任务4"}
	keeper.AddJob(3*time.Second, job4)

	job := &JobTest{name: "任务1"}
	keeper.AddJob(2*time.Second, job)
	job2 := &JobTest{name: "任务2"}
	keeper.AddJob(2*time.Second, job2)

	job5 := &JobTest{name: "任务5"}
	keeper.AddJob(0, job5)
	keeper.ExecJob()
	select {}
}
