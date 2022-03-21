package generic3

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Method interface {
	Init(Module) error
	Wait(ctx context.Context) error
	Done()
	Do()
}

type SyncMethod Method
type AsyncMethod Method

type Module interface {
	AddTask(req Method) error
}

type BaseMethod[T Module] struct {
	module T
}

func (b *BaseMethod[T]) Init(t Module) error {
	var m, ok = t.(T)
	if !ok {

	}
	b.module = m
	return nil
}
func (b *BaseMethod[T]) Wait(ctx context.Context) error { return nil }
func (b *BaseMethod[T]) Done()                          {}
func (b *BaseMethod[T]) GetModule() T                   { return b.module }

type BaseSyncMethod[M Module, T any] struct {
	BaseMethod[M]
	done chan struct{}
	rsp  T
}

func (b *BaseSyncMethod[M, T]) Init(t Module) error {
	if err := b.BaseMethod.Init(t); err != nil {
		return err
	}
	b.done = make(chan struct{})
	return nil
}

func (b *BaseSyncMethod[M, T]) Wait(ctx context.Context) error {
	select {
	case <-b.done:
		return nil
	case <-ctx.Done():
		return ErrWaitTimeout
	}
}

func (b *BaseSyncMethod[M, T]) Done() {
	close(b.done)
}

func (b *BaseSyncMethod[M, T]) Resp() *T {
	return &b.rsp
}

type BaseAsyncMethod[M Module] struct {
	BaseMethod[M]
}

type MyModule struct {
	name        string
	worker      Module // 实际工作的模块
	closeChan   chan struct{}
	commendChan chan Method
}

var (
	ErrAddTimeout  = fmt.Errorf("add task timeout")
	ErrWaitTimeout = fmt.Errorf("wait task timeout")
)

func (m *MyModule) AddTask(req Method) error {
	if err := req.Init(m.worker); err != nil {
		return err
	}
	var ctx, cancel = context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	select {
	case m.commendChan <- req:
	case <-ctx.Done():
		return ErrAddTimeout
	}
	if err := req.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func (m *MyModule) Run() {

	m.closeChan = make(chan struct{})
	m.commendChan = make(chan Method, 128)

	go func() {
		for {
			select {
			case <-m.closeChan:
				return
			case cmd := <-m.commendChan:
				cmd.Do()
				cmd.Done()
			}
		}
	}()
}

func (m *MyModule) Init() {
	m.worker = m
}

type MyReq struct {
	BaseSyncMethod[*MyModule, MyResp]
	Name string
}

func (m *MyReq) Do() {
	var rsp = m.Resp()
	rsp.Name = m.Name + m.GetModule().name
	rsp.Age = 100
}

type MyResp struct {
	Name string
	Age  int
}

type MyAsyncReq struct {
	BaseAsyncMethod[*MyModule]
	Haha string
}

func (m *MyAsyncReq) Do() {
	var mm = m.GetModule()
	log.Println(mm.name + " " + m.Haha)

}
