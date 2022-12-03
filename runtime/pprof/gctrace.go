package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func alloc() int {
	var m = make([]int, 1<<20)
	runtime.SetFinalizer(&m, func(obj interface{}) {
		log.Println("Finalizer m")
	})
	runtime.KeepAlive(m)
	return 0
}

// 内存: go tool pprof -alloc_space http://localhost:9999/debug/pprof/heap
//

func main() {
	someError()
}

func allocProfile() {
	var base, _ = os.Getwd()
	log.Println(base)

	var cpu CPUProfile
	var mem MemProfile
	var h HTTPProfile

	defer Run(&cpu, base+"/cpu.pprof")()
	defer Run(&mem, base+"/mem.pprof")()
	defer Run(&h, "localhost:8899")()

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	var tick = time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			var ret = alloc()
			log.Println(ret)
		case <-sigChan:
			log.Println("[WRN] stop")
			return
		}
	}
}
