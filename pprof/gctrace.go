package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
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

	var cpu CPUProfile
	var mem MemProfile
	var h HTTPProfile

	defer Run(&cpu, base+"/cpu.pprof")()
	defer Run(&mem, base+"/mem.pprof")()
	defer Run(&h, "localhost:8899")()

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		for i := 0; i < 10; i++ {
			var ret = alloc()
			log.Println(ret)
		}
	}()
	go func() {
		_ = http.ListenAndServe(":9999", nil)
	}()

	wait.Wait()
}
