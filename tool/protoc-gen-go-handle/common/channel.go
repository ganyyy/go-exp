package common

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	ErrModuleClosed   = errors.New("module closed")
	ErrTimeout        = errors.New("timeout")
	ErrNotInitialized = errors.New("not initialized")
)

type Task struct {
	Context context.Context
	Code    string
	Req     proto.Message
	Rsp     proto.Message
	Wait    chan error
}

func (t Task) Finish(err error) {
	if t.Wait != nil {
		select {
		case t.Wait <- err:
		default:
		}
	}
}

type ChannelInvoke struct {
	taskChan chan Task
	done     context.Context
	stop     context.CancelFunc
}

// isInitialized 是否初始化
func (ch *ChannelInvoke) isInitialized() bool {
	return ch.taskChan != nil
}

// Init 初始化
func (ch *ChannelInvoke) Init(ctx context.Context, capacity uint) {
	ch.taskChan = make(chan Task, capacity)
	ch.done, ch.stop = context.WithCancel(ctx)
}

// Done 获取完成上下文
func (ch *ChannelInvoke) Done() <-chan struct{} {
	if !ch.isInitialized() {
		panic(ErrNotInitialized)
	}
	return ch.done.Done()
}

// TaskChan 获取任务通道
func (ch *ChannelInvoke) TaskChan() <-chan Task {
	if !ch.isInitialized() {
		panic(ErrNotInitialized)
	}
	return ch.taskChan
}

// Stop 停止
func (ch *ChannelInvoke) Stop() {
	if !ch.isInitialized() {
		panic(ErrNotInitialized)
	}
	ch.stop()
}

// Invoke 处理请求
func (ch *ChannelInvoke) Invoke(ctx context.Context, code string, req proto.Message, rsp proto.Message) error {
	if !ch.isInitialized() {
		return ErrNotInitialized
	}
	var wait chan error
	if rsp != nil { // 同步请求, 存在返回值
		wait = make(chan error, 1)
	}
	task := Task{
		Context: ctx,
		Code:    code,
		Req:     req,
		Rsp:     rsp,
		Wait:    wait,
	}
	select {
	case ch.taskChan <- task:
	case <-ctx.Done():
		return ctx.Err()
	case <-ch.Done():
		return ErrModuleClosed
	}
	if wait == nil {
		return nil // 异步
	}

	select { // 同步
	case err := <-wait:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-ch.Done():
		return ErrModuleClosed
	}
}
