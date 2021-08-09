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

func main() {

	var useRaw = flag.Bool("r", false, "use raw syscall")
	var timeSleep = flag.Int("t", 1000, "sleep ms time")
	var loopNum = flag.Int("n", 3, "loop time")

	flag.Parse()
	var ms = *timeSleep

	log.Printf("use raw syscall:%v, time wait:%v, loop num:%v", *useRaw, ms, *loopNum)
	time.Sleep(time.Second * 2)

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGURG)
	go func() {
		for sig := range sigChan {
			log.Printf("[%v] recv sig:%v", time.Now().UnixNano()/int64(time.Millisecond), sig)
		}
	}()

	var selectCall func() error
	var nfd = syscall.Stderr + 1

	if *useRaw {
		selectCall = func() error {
			var timeout = &syscall.Timeval{
				Usec: 1000 * int64(ms),
			}
			_, _, e1 := syscall.RawSyscall6(
				syscall.SYS_SELECT,
				uintptr(nfd), uintptr(0),
				uintptr(0), uintptr(0),
				uintptr(unsafe.Pointer(timeout)), 0,
			)
			return e1
		}
	} else {
		selectCall = func() error {
			var timeout = &syscall.Timeval{
				Usec: 1000 * int64(ms),
			}
			_, err := syscall.Select(nfd, nil, nil, nil, timeout)
			return err
		}
	}

	var cnt int
	for {
		e1 := selectCall()
		if e1 == nil || e1 == syscall.Errno(0) || e1 == syscall.EINTR {
			log.Printf("[%v] retry, err:%v", time.Now().UnixNano()/int64(time.Millisecond), e1)
			cnt++
			if cnt > *loopNum {
				break
			}
		} else {
			log.Fatalf("other error:%v", e1)
		}
	}
}
