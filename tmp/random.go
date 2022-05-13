// +build ignore

package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	var wait sync.WaitGroup

	var cpuNum = runtime.GOMAXPROCS(0)

	wait.Add(cpuNum)

	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	var cnt uint64
	func() {
		for i := 0; i < cpuNum; i++ {
			go func(i int) {
				defer func() {
					if recover() != nil {
						log.Printf("%v total rand num:%v", i, atomic.LoadUint64(&cnt))
					}
				}()
				defer wait.Done()
				for i := 0; i < (1 << 30); i++ {
					atomic.AddUint64(&cnt, 1)
					_ = rng.Int63()
				}
			}(i)
		}
		wait.Wait()
	}()

}
