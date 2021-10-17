//go:build linux
// +build linux

package main

import (
	"go-exp/helper"
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	var err error

	// epoll fd
	var epfd int
	var events [5]unix.EpollEvent
	epfd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	helper.PanicIfErr("create epollfd error", err)

	// event fd
	var efd int
	var eval uint64 // eventfd 读写缓冲区大小是固定的, 只有8个字节
	var ebuff = (*(*[8]byte)(unsafe.Pointer(&eval)))[:]
	efd, _ = unix.Eventfd(0, 0)
	unix.SetNonblock(efd, true)

	unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, efd, &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLET,
		Fd:     int32(efd),
	})

	// timer fd
	var timer unix.ItimerSpec
	var timerfd int
	timer.Value.Sec = 1    // 初始1秒定时
	timer.Interval.Sec = 1 // 间隔1秒一次
	timerfd, _ = unix.TimerfdCreate(unix.CLOCK_REALTIME, 0)
	unix.TimerfdSettime(timerfd, unix.TFD_TIMER_ABSTIME, &timer, nil)

	unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, timerfd, &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLET,
		Fd:     int32(timerfd),
	})

	go func() {
		for {
			var n, err = unix.EpollWait(epfd, events[:], -1)
			if n < 0 {
				if err == unix.EAGAIN || err == unix.EINTR {
					continue
				}
				log.Fatalf("epoll wait error:%v", err)
			}
			for _, e := range events[:n] {
				if e.Fd == int32(efd) {
					var val uint64
					unix.Read(efd, (*(*[8]byte)(unsafe.Pointer(&val)))[:])
					log.Printf("efd:%v, val:%v", efd, val)
				} else if e.Fd == int32(timerfd) {
					var val uint64
					_, _ = unix.Read(timerfd, (*(*[8]byte)(unsafe.Pointer(&val)))[:])
					log.Printf("timerfd:%v, val:%v", timerfd, val)
				}
			}
		}
	}()

	{
		// 写event fd
		// event fd中的值是累加的, 也就是说, 读取到的是累加值(200+300)
		eval = 200
		unix.Write(efd, ebuff)
		time.Sleep(time.Second)
		eval = 200
		unix.Write(efd, ebuff)
	}

	{

	}

	time.Sleep(time.Minute * 5)
}
