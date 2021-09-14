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

	var call func() error
	var nfd int
	nfd, _ = syscall.Kqueue()

	//var fd, _ = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	var events [1]syscall.Kevent_t

	if *useRaw {
		call = func() error {
			var timeout = &syscall.Timespec{
				Nsec: int64(1000 * 1000 * (*timeSleep)),
			}
			_, _, e1 := syscall.RawSyscall6(
				syscall.SYS_KEVENT,
				uintptr(nfd), 0, 0,
				uintptr(unsafe.Pointer(&events[0])), uintptr(len(events)),
				uintptr(unsafe.Pointer(timeout)),
			)
			return e1
		}
	} else {
		call = func() error {
			var timeout = &syscall.Timespec{
				Nsec: int64(1000 * 1000 * (*timeSleep)),
			}
			_, err := syscall.Kevent(nfd, nil, events[:], timeout)
			return err
		}
	}

	var cnt int
	for {
		e1 := call()
		if e1 == nil || e1 == syscall.Errno(0) || e1 == syscall.EINTR {
			log.Printf("[%v] retry, err:%v", time.Now().UnixNano()/int64(time.Millisecond), e1)
			cnt++
			if cnt >= *loopNum {
				break
			}
		} else {
			log.Fatalf("other error:%v", e1)
		}
	}
}
