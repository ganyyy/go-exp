package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

var wait sync.WaitGroup
var mutex sync.Mutex
var dataChan = make(chan int)

func do(ctx context.Context, f func(ctx context.Context)) {
	if f == nil {
		return
	}
	wait.Add(1)
	go func() {
		defer wait.Done()
		f(ctx)
	}()
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
		ticker := time.NewTicker(time.Second)
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
		case dataChan <- 10:
		}
	}
}

func sendBlock(ctx context.Context) {
	_ = ctx
	dataChan <- 1
}

func someError() {
	var ctx, cancel = context.WithCancel(context.Background())

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	do(ctx, mutexBlock)
	do(ctx, mutexBlock)
	do(ctx, ticker)
	do(ctx, slice)
	do(ctx, chanBlock)

	go sendBlock(ctx)

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
