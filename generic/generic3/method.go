package generic3

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type ErrorMsg interface {
	Code() Code
	CheckCode(code Code) bool
	SetCode(code Code) bool
	SetError(code Code, err interface{}) bool
	Error() error
}

type Method interface {
	ErrorMsg

	Init(Module) error
	Do()
	Done()
}

type Code = uint32

const (
	CodeNo Code = iota
	CodeAddTimeout
	CodeWaitTimeout
	CodeModuleNotMatch
	CodeInner
)

type SyncMethod Method
type AsyncMethod Method

type Module interface {
	AddTask(req Method)
}

type Base struct {
	msg  interface{}
	code Code
}

func (b *Base) Code() Code {
	return Code(atomic.LoadUint32((*uint32)(&b.code)))
}

func (b *Base) CheckCode(code Code) bool {
	return b.Code() == code
}

func (b *Base) SetCode(code Code) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&b.code), uint32(CodeNo), uint32(code))
}

func (b *Base) SetError(code Code, err interface{}) bool {
	if !b.SetCode(code) {
		return false
	}
	b.msg = err
	return true
}

func (b *Base) Error() error {
	if b.CheckCode(CodeNo) {
		return nil
	}
	return fmt.Errorf("<code %v, err %v>", b.code, b.msg)
}

type BaseMethod[T Module] struct {
	Base
	module T
}

func (b *BaseMethod[T]) Init(t Module) error {
	var m, ok = t.(T)
	if !ok {
		b.SetCode(CodeModuleNotMatch)
		b.Done()
	}
	b.module = m
	return nil
}
func (b *BaseMethod[T]) GetModule() T { return b.module }
func (b *BaseMethod[T]) Done()        {}

type BaseSyncMethod[M Module, T any] struct {
	BaseMethod[M]
	rsp     T
	rspChan chan T
}

func (b *BaseSyncMethod[M, T]) Init(t Module) error {
	if err := b.BaseMethod.Init(t); err != nil {
		return err
	}
	b.rspChan = make(chan T, 1)
	return nil
}

func (b *BaseSyncMethod[M, T]) GenResp() *T {
	return &b.rsp
}

func (b *BaseSyncMethod[M, T]) SetResp(resp T) {
	select {
	case b.rspChan <- resp:
	default:
	}
}

func (b *BaseSyncMethod[M, T]) Resp() (T, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return b.RespWithContext(ctx)
}

func (b *BaseSyncMethod[M, T]) RespWithContext(ctx context.Context) (T, error) {
	select {
	case rsp := <-b.rspChan:
		return rsp, b.Error()
	case <-ctx.Done():
		var empty T
		return empty, ErrWaitTimeout
	}
}

func (b *BaseSyncMethod[M, T]) Done() {
	b.SetResp(b.rsp)
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

func (m *MyModule) AddTask(req Method) {
	if err := req.Init(m.worker); err != nil {
		req.SetCode(CodeModuleNotMatch)
		req.Done()
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case m.commendChan <- req:
	case <-ctx.Done():
		req.SetError(CodeAddTimeout, ErrAddTimeout)
		req.Done()
	}
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
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println(string(debug.Stack()))
							cmd.SetError(CodeInner, err)
						}
						cmd.Done()
					}()
					log.Printf("%+v", cmd)
					if !cmd.CheckCode(CodeNo) {
						return
					}
					cmd.Do()
				}()

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
	var rsp = m.GenResp()
	rsp.Name = m.Name + m.GetModule().name
	rsp.Age = 100

}

type MyResp struct {
	Base
	Name string
	Age  int
}

type MyAsyncReq struct {
	BaseAsyncMethod[*MyModule]
	Lala string
}

func (m *MyAsyncReq) Do() {
	var mm = m.GetModule()
	log.Println(mm.name + " " + m.Lala)

}
