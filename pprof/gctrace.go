package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func alloc() int {
	var m = make([]int, 1<<20)
	runtime.SetFinalizer(&m, func(obj interface{}) {
		log.Println("Finalizer m")
	})
	time.Sleep(time.Minute)
	return m[0]
}

func main() {
	go func() {
		var ret = alloc()
		log.Println(ret)
	}()
	_ = http.ListenAndServe(":9999", nil)
}
