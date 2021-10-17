//go:build darwin

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

//go:generate go build main.go

/*
export GODEBUG=asyncpreemptoff=1

unset GODEBUG
*/

func main() {

	var useRaw = flag.Bool("r", false, "use raw syscall")
	var timeSleep = flag.Int("t", 1, "sleep ms time")
	var loopNum = flag.Int("n", 3, "loop time")

	flag.Parse()
	var ms = *timeSleep

	log.Printf("use raw syscall:%v, time wait:%v, loop num:%v", *useRaw, ms, *loopNum)
	time.Sleep(time.Second * 2)

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGURG)
	go func() {
		for sig := range sigChan {
			log.Printf("[%v] recv sig:%v", time.Now().UnixMilli(), sig)
		}
	}()

	var sysCall func() error

	if *useRaw {
		sysCall = func() error {
			var timeout = &syscall.Timespec{
				Sec: int64(ms),
			}
			_, _, e1 := syscall.RawSyscall(
				syscall.SYS_NANOSLEEP,
				uintptr(unsafe.Pointer(timeout)), 0, 0)
			return e1
		}
	} else {
		sysCall = func() error {
			var timeout = &syscall.Timespec{
				Sec: int64(ms),
			}
			err := syscall.Nanosleep(timeout, nil)
			return err
		}
	}

	var cnt int
	for {
		e1 := sysCall()
		if e1 == nil || e1 == syscall.Errno(0) || e1 == syscall.EINTR {
			log.Printf("[%v] retry, err:%v", time.Now().UnixMilli(), e1)
			cnt++
			if cnt >= *loopNum {
				break
			}
		} else {
			log.Fatalf("other error:%v", e1)
		}
	}
}
