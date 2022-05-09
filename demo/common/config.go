package common

import (
	"errors"
	"time"
)

var (
	errChanSize    = errors.New("invalid task chan size")
	errEmptyWorker = errors.New("must assign valid worker module")
	errTimeout     = errors.New("must input valid timeout")
)

// worker的配置信息
type workerConfig struct {
	logPrefix string        // 日志前缀
	chanSize  int           // 任务队列的长度
	worker    IWorker       // 实际的工作单元, 必填项目
	isDebug   bool          // 开启debug模式, 会输出详细的任务日志, 线上请关闭这个
	waitAsync bool          // 是否需要等待异步任务进入队列
	timeout   time.Duration // 等待时常

	//TODO 添加新的配置一定要提供一个初始化的值
	//     以及一个赋值Option函数
	//	   必要的话, 在check中添加检查函数
}

func (w workerConfig) check() error {
	if w.chanSize <= 0 {
		return errChanSize
	}
	if w.worker == nil {
		return errEmptyWorker
	}
	// 至少要等个1s吧
	if w.timeout == 0 || w.timeout < time.Second {
		return errTimeout
	}
	return nil
}

//NewWorker 创建一个新的模块, 子类需要内嵌 IWorker, 再将自己传入进来. 并将返回的结果进行替换
func NewWorker(worker IWorker, configs ...WorkerConfigOpt) (IWorker, error) {
	var config = workerConfig{
		logPrefix: "[Worker]",          // 默认的日志前缀
		chanSize:  DefaultWorkChanSize, // 默认的task chan大小
		worker:    worker,              // 实际工作模块
		isDebug:   false,               // 是否开启Debug模式
		waitAsync: true,                // 默认等待异步任务进入队列
		timeout:   5 * time.Second,     // 默认等待超时的时常
	}

	for _, cfg := range configs {
		cfg(&config)
	}

	if err := config.check(); err != nil {
		return nil, err
	}

	return &baseWorker{workerConfig: config}, nil
}

type WorkerConfigOpt func(*workerConfig)

func ConfigPrefix(logPrefix string) WorkerConfigOpt {
	return func(config *workerConfig) {
		config.logPrefix = logPrefix
	}
}

func ConfigChanSize(size int) WorkerConfigOpt {
	return func(config *workerConfig) {
		config.chanSize = size
	}
}

func ConfigDebugMode(isDebug bool) WorkerConfigOpt {
	return func(config *workerConfig) {
		config.isDebug = isDebug
	}
}

func ConfigTimeout(timeout time.Duration) WorkerConfigOpt {
	return func(config *workerConfig) {
		config.timeout = timeout
	}
}
