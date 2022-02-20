package main

import (
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	var fd, _ = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)

	unix.SetNonblock(fd, true)

	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)

	unix.Bind(fd, &unix.SockaddrInet4{
		Port: 6666,
		Addr: [4]byte{
			127,
			0,
			0,
			1,
		},
	})

	go func() {
		unix.Connect(fd, &unix.SockaddrInet4{
			Port: 6666,
			Addr: [4]byte{
				127,
				0,
				0,
				1,
			},
		})
	}()

	time.Sleep(time.Minute)
}
