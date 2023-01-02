package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var wait sync.WaitGroup
var mutex sync.Mutex
var dataSendChan = make(chan int)
var dataRecvChan = make(chan int)
var deadlock sync.Mutex

func do(ctx context.Context, name string, f func(ctx context.Context)) {
	if f == nil {
		return
	}
	wait.Add(1)
	go pprof.Do(ctx, pprof.Labels("name", name), func(ctx context.Context) {
		defer wait.Done()
		f(ctx)
	})
}

func slice(ctx context.Context) {
	var data [][]byte
end:
	for {
		select {
		case <-ctx.Done():
			break end
		default:
		}
		data = append(data, allocBytes())
		time.Sleep(time.Second)
	}
	log.Printf("all data:%v", len(data))
}

func mutexBlock(ctx context.Context) {
end:
	for {
		select {
		case <-ctx.Done():
			break end
		default:
		}
		mutex.Lock()
		time.Sleep(time.Second * 2)
		mutex.Unlock()
	}
	log.Printf("check done")
}

func ticker(ctx context.Context) {
	for {
		ticker := time.NewTicker(time.Millisecond * 100)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// dataChan <- int(t.Unix())
		}
	}
}

func chanBlock(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case dataSendChan <- 10:
		}
	}
}

func groupWait(ctx context.Context) {
	var wait sync.WaitGroup
	wait.Add(1)
	wait.Wait()
}

func sendBlock(ctx context.Context) {
	_ = ctx
	dataSendChan <- 1
}

func recvBlock(ctx context.Context) {
	_ = ctx
	<-dataRecvChan
}

func deadLock(ctx context.Context) {
	deadlock.Lock()
}

func someError() {

	log.Printf("Default rate:%v", runtime.MemProfileRate)

	// 可以适当的降低一下采样的频率
	const RATE = 16 << 10
	runtime.MemProfileRate = RATE

	var ctx, cancel = context.WithCancel(context.Background())

	// 开启竞态检查
	runtime.SetMutexProfileFraction(1)
	// 开启阻塞检查
	runtime.SetBlockProfileRate(1)

	do(ctx, "mutexBlock", mutexBlock)
	do(ctx, "mutexBlock", mutexBlock)
	do(ctx, "ticker", ticker)
	do(ctx, "slice", slice)
	do(ctx, "chanBlock", chanBlock)

	go groupWait(ctx)
	go sendBlock(ctx)
	go recvBlock(ctx)
	deadlock.Lock()
	go deadLock(ctx)
	go func() {
		err := http.ListenAndServe(":9999", nil)
		if err != nil {
			log.Printf("listen error:%v", err)
			os.Exit(1)
		}
	}()

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	log.Printf("recv %v", sig)
	cancel()
	wait.Wait()
	log.Printf("done")
}

func allocBytes() []byte {
	return make([]byte, 1<<10)[:5]
}
