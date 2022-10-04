package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

func deadBlock() {
	var mutex sync.Mutex
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	var check = func() {
		for {
			mutex.Lock()
			time.Sleep(time.Second * 20)
			mutex.Unlock()
		}
	}
	go check()
	go check()

	_ = http.ListenAndServe(":9999", nil)
}
