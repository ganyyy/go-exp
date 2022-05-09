package common

import (
	"context"
)

type Handle func()

//ITask 纯内部接口, 外部需要通过继承 SyncTaskBase, AsyncTaskBase 的形式实现
type ITask interface {
	failed(innerErrCode, interface{}) // 内部错误
	init()                            // 任务信息初始化
	wait(context.Context)             // 等待任务执行完成
	sync() bool                       // 是不是同步的任务
	finish()                          // 任务结束通知, 内部调用
	error() error                     // 内部错误信息
}

//IErrorCode 错误吗接口, 自己实现一个就行. 推荐使用 stringer生成
type IErrorCode interface {
	Code() uint32
	String() string
	Error() string
	Equal(IErrorCode) bool
}

//IWorker 工作模块结构, 都是导出的方法是为了方便外部重写
type IWorker interface {
	Init()
	Start()
	Stop()
	LogName() string

	AddTask(ITask) error   // 添加任务到队列中, 如果有错, 则结果不可用, 这个方法千万不要重写
	HandleTask(ITask)      // 处理任务, 这个方法千万不要重写
	Task() <-chan ITask    // 任务队列, 这个方法千万不要重写
	Done() <-chan struct{} // 结束通知, 这个方法千万不要重写

	//TODO 以下是要重写的函数, 就重写这三个就行了

	CheckTask(ITask) bool // need override 检查任务的类型
	DoTask(ITask)         // need override 要重写的函数
	Handle()              // optional override 逻辑处理函数, 如果没有定时任务啥的可以不重写
}
