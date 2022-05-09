package common

import (
	"golang.org/x/net/context"
)

type taskBase struct {
	info ErrorInfo
}

func (t *taskBase) init()                                   {}
func (t *taskBase) wait(_ context.Context)                  {}
func (t *taskBase) finish()                                 {}
func (t *taskBase) sync() bool                              { return false }
func (t *taskBase) error() error                            { return t.info.Error() }
func (t *taskBase) failed(code innerErrCode, i interface{}) { t.info.Set(code, i) }

//AsyncTaskBase 异步任务
type AsyncTaskBase struct {
	taskBase
}

//SyncTaskBase 同步任务
type SyncTaskBase struct {
	taskBase
	done chan struct{}
}

func (t *SyncTaskBase) init() {
	t.done = make(chan struct{})
}

func (t *SyncTaskBase) wait(ctx context.Context) {
	select {
	case <-t.done:
	case <-ctx.Done():
		t.failed(ErrWaitTimeout, nil)
	}
}

func (t *SyncTaskBase) finish() {
	close(t.done)
}

func (t *SyncTaskBase) sync() bool { return true }
