package common_test

import (
	"log"
	"testing"

	"ganyyy.com/go-exp/demo/common"
	"github.com/stretchr/testify/assert"
)

type TestWorker struct {
	common.IWorker // 实际的执行单元
	value          int
}

type TestTask interface {
	common.ITask
	Do(*TestWorker)
}

type TestTask1 struct {
	common.SyncTaskBase

	Resp struct {
		V int
	}
}

func (t *TestTask1) Do(worker *TestWorker) {
	t.Resp.V = worker.value
}

type TestTask2 struct {
	common.AsyncTaskBase
	V int
}

func (t *TestTask2) Do(worker *TestWorker) {
	worker.value = t.V
}

func (t *TestWorker) CheckTask(task common.ITask) bool {
	_, ok := task.(TestTask)
	return ok
}

func (t *TestWorker) DoTask(task common.ITask) {
	task.(TestTask).Do(t)
}

func (t *TestWorker) Handle() {
	log.Printf("In Outer handle")
	for {
		select {
		case task := <-t.Task():
			log.Printf("task info:%v", common.TaskToDetail(task))
			t.HandleTask(task)
		case <-t.Done():
			return
		}
	}
}

func TestWorkerOutFunc(t *testing.T) {
	var testWorker TestWorker
	var worker, err = common.NewWorker(
		&testWorker,
		common.ConfigPrefix("[12345]"),
		common.ConfigDebugMode(true),
	)
	assert.Nil(t, err)
	testWorker.IWorker = worker
	testWorker.Init()
	testWorker.Start()
	defer testWorker.Stop()
	err = testWorker.AddTask(&TestTask2{
		V: 200,
	})
	assert.Nil(t, err)

	var t1 TestTask1
	err = testWorker.AddTask(&t1)
	assert.Nil(t, err)
	assert.Equal(t, t1.Resp.V, 200)
}
