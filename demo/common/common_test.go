package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestWorker struct {
	IWorker // 实际的执行单元
	value   int
}

type TestTask interface {
	ITask
	Do(*TestWorker)
}

type TestTask1 struct {
	SyncTaskBase

	Resp struct {
		V int
	}
}

func (t *TestTask1) Do(worker *TestWorker) {
	t.Resp.V = worker.value
}

type TestTask2 struct {
	AsyncTaskBase
	V int
}

func (t *TestTask2) Do(worker *TestWorker) {
	worker.value = t.V
}

func (t *TestWorker) CheckTask(task ITask) bool {
	_, ok := task.(TestTask)
	return ok
}

func (t *TestWorker) DoTask(task ITask) {
	task.(TestTask).Do(t)
}

func TestWorkerFunc(t *testing.T) {
	var testWorker TestWorker
	var worker, err = NewWorker(&testWorker)
	assert.Nil(t, err)
	testWorker.IWorker = worker
	testWorker.Init()
	testWorker.Start()
	defer testWorker.Stop()
	err = testWorker.AddTask(&TestTask2{
		V: 100,
	})
	assert.Nil(t, err)

	var t1 TestTask1
	err = testWorker.AddTask(&t1)
	assert.Nil(t, err)
	assert.Equal(t, t1.Resp.V, 100)
}
