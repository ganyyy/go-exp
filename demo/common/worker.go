package common

import (
	"context"
	"log"
	"sync"
	"time"
)

//baseWorker 基础实现, 不需要暴露给外部
type baseWorker struct {
	wg           sync.WaitGroup
	workerConfig               // 相关的配置
	closeChan    chan struct{} // 关闭队列
	taskChan     chan ITask    // 执行任务队列

	startOnce sync.Once
	closeOnce sync.Once
}

func (b *baseWorker) Handle() {
	for {
		select {
		case task := <-b.Task():
			b.HandleTask(task)
		case <-b.Done():
			return
		}
	}
}

func (b *baseWorker) HandleTask(task ITask) {
	defer func() {
		if err := recover(); err != nil {
			task.failed(ErrInner, err)
			log.Printf("%v handle task %v crash, err: %v", b.LogName(), TaskName(task), err)
		}
		task.finish()
	}()
	if b.isDebug {
		defer func(begin time.Time) {
			log.Printf("%v Do [%v][%v] %v",
				b.LogName(), TaskName(task), time.Since(begin), TaskToDetail(task),
			)
		}(time.Now())
	}
	// 父类调用子类的方法, 学废了吗
	b.worker.DoTask(task)
}

func (b *baseWorker) Task() <-chan ITask {
	return b.taskChan
}

func (b *baseWorker) Done() <-chan struct{} {
	return b.closeChan
}

func (b *baseWorker) Init() {
	b.closeChan = make(chan struct{})
	b.taskChan = make(chan ITask, b.chanSize)
}

func (b *baseWorker) Start() {
	b.startOnce.Do(func() {
		log.Printf("%v start!", b.LogName())
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			b.worker.Handle()
		}()
	})
}

func (b *baseWorker) Stop() {
	b.closeOnce.Do(func() {
		close(b.closeChan)
		b.wg.Wait()
		log.Printf("%v stop!", b.LogName())
	})
}

func (b *baseWorker) LogName() string {
	return b.logPrefix
}

func (b *baseWorker) AddTask(task ITask) error {
	if !b.worker.CheckTask(task) {
		task.failed(ErrTaskType, nil)
		return task.error()
	}
	task.init()
	if task.sync() || b.waitAsync {
		// 如果是同步任务; 或者配置为等待异步任务进入队列(不保证可以完成)
		var ctx, cancel = context.WithTimeout(context.Background(), b.timeout)
		defer cancel()

		select {
		case b.taskChan <- task:
		case <-ctx.Done():
			task.failed(ErrAddTimeout, nil)
			return task.error()
		}
		task.wait(ctx)
	} else {
		// 如果不关心异步任务是否执行完成, 直接添加不进去直接返回即可
		select {
		case b.taskChan <- task:
		default:
			task.failed(ErrChanFull, nil)
		}
	}
	return task.error()
}

// 以下两个函数, 应该由子类重写

func (b *baseWorker) CheckTask(_ ITask) bool {
	return false
}

func (b *baseWorker) DoTask(task ITask) {
	task.failed(ErrInner, "not implement")
}
