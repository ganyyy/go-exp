//go:build ignore

package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	//MultiCond()

	MultiAtomic()
}

func MultiGo() {
	const (
		Worker  = 10
		Timeout = 5
	)

	var ctx, cancel = context.WithCancel(context.TODO())
	var wait sync.WaitGroup

	wait.Add(Worker)
	for i := 0; i < Worker; i++ {
		go func(id int) {
			defer wait.Done()
			var tick = time.NewTicker(time.Second)
			var cnt int
			for {
				select {
				case <-tick.C:
					log.Printf("[%v]Tick Run %v", id, cnt)
					cnt++
				case <-ctx.Done():
					log.Printf("[%v]Exit", id)
					return
				}
			}
		}(i)
	}

	<-time.After(time.Second*Timeout + time.Millisecond*100)
	cancel()

	wait.Wait()

	log.Printf("Main Goroutine Done!")
}

func MultiCond() {

	const (
		Worker = 10
	)

	var mutex sync.Mutex
	var cond = sync.NewCond(&mutex)
	var wait sync.WaitGroup

	wait.Add(Worker)

	var val uint32

	for i := 0; i < Worker; i++ {
		go func(id int) {
			defer wait.Done()
			mutex.Lock()
			for val != uint32(id) {
				cond.Wait()
			}
			var old = val
			val++
			//time.Sleep(time.Millisecond * 100)
			cond.Broadcast()
			mutex.Unlock()
			log.Printf("Goroutine:%v, Val:%v", id, old)
		}(i)
	}

	wait.Wait()
}

func MultiAtomic() {
	const (
		Worker = 100
	)

	var wait sync.WaitGroup
	wait.Add(Worker)
	defer wait.Wait()

	var val uint32

	for i := 0; i < Worker; i++ {
		go func(id int) {
			defer wait.Done()
			for atomic.LoadUint32(&val) != uint32(id) {
				runtime.Gosched()
			}
			log.Printf("Goroutine:%v, val:%v", id, val)
			atomic.AddUint32(&val, 1)
		}(i)
	}
}
