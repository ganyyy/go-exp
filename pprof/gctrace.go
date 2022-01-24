package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func alloc() int {
	var m = make([]int, 1<<20)
	runtime.SetFinalizer(&m, func(obj interface{}) {
		log.Println("Finalizer m")
	})
	time.Sleep(time.Second)
	runtime.KeepAlive(m)
	return 0
}

// 内存: go tool pprof -alloc_space http://localhost:9999/debug/pprof/heap
//

func main() {

	var base, _ = os.Getwd()
	log.Println(base)

	defer (&CPUProfile{}).Init(base + "/cpu.pprof").Start().Done()
	defer (&MemProfile{}).Init(base + "/mem.pprof").Start().Done()
	defer (&HTTPProfile{}).Init("localhost:9999").Start().Done()

	go func() {
		for i := 0; i < 10; i++ {
			var ret = alloc()
			log.Println(ret)
		}
	}()
	go func() {
		_ = http.ListenAndServe(":9999", nil)
	}()

	time.Sleep(time.Minute * 5)
}
