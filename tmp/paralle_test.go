package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

type readParam struct {
	m    map[int]int
	done chan bool
}

const (
	randNum = 10000
)

func runWorker() (chan<- readParam, chan<- bool) {
	var worker = make(chan readParam, 1024)
	var done = make(chan bool, 1)

	var src = make(map[int]int, randNum)

	for i := 0; i < randNum; i++ {
		src[i] = rand.Intn(randNum)
	}

	go func() {
		for {
			select {
			case cmd := <-worker:
				for id := range cmd.m {
					cmd.m[id] = src[id]
				}
				close(cmd.done)
			case <-done:
				return
			}
		}
	}()

	return worker, done
}

func TestParallelReadMap(t *testing.T) {
	var worker, done = runWorker()

	var wait sync.WaitGroup

	const (
		WorkerNum = 10
	)

	var getRet = func(done chan struct{}) (bool, map[int]int) {
		var cmd readParam
		cmd.m = make(map[int]int, 100)
		cmd.done = make(chan bool)

		for i := 0; i < 100; i++ {
			cmd.m[rand.Intn(randNum)] = 0
		}

		select {
		case worker <- cmd:
		case <-done:
			return false, cmd.m
		}

		select {
		case <-cmd.done:
			return true, cmd.m
		case <-done:
			return false, cmd.m
		}
	}

	wait.Add(WorkerNum)

	for i := 0; i < WorkerNum; i++ {
		go func(work int) {
			defer wait.Done()
			defer log.Printf("work %v done", work)
			for c := 0; c < 10; c++ {
				var dd = make(chan struct{})
				go func(w, c int) {
					var ret, m = getRet(dd)
					if ret {
						var r int
						for _, v := range m {
							r += v
						}
						log.Printf("%v:%v Ret %v Sum %v", w, c, ret, r)
					}

				}(work, c)
				runtime.Gosched()
				close(dd)
			}
		}(i)
	}

	wait.Wait()

	close(done)
}
