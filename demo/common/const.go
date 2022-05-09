package common

const (
	ErrNo          innerErrCode = iota // 没毛病
	ErrWaitTimeout                     // 等待执行结束超时
	ErrAddTimeout                      // 添加任务到队列超时
	ErrInner                           // 内部错误
	ErrTaskType                        // 任务类型错误
	ErrChanFull                        // 队列满了
)

const (
	DefaultWorkChanSize = 512
)
