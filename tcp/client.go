//go:build ignore
// +build ignore

package main

import (
	"syscall"
	"time"

	"go-exp/helper"
)

func main() {

	var fd, fd2 int
	var err error

	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	helper.PanicIfErr("create socket", err)
	fd2, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	helper.PanicIfErr("create socket", err)

	syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: 55555,
		Addr: [4]byte{
			127, 0, 0, 1,
		},
	})

	syscall.Bind(fd2, &syscall.SockaddrInet4{
		Port: 55555,
		Addr: [4]byte{
			172, 17, 198, 197,
		},
	})

	var serAddr = &syscall.SockaddrInet4{
		Port: 9999,
		Addr: [4]byte{
			127, 0, 0, 1,
		},
	}

	err = syscall.Connect(fd, serAddr)
	helper.PanicIfErr("fd1 connect", err)
	err = syscall.Connect(fd2, serAddr)
	helper.PanicIfErr("fd2 connect", err)

	time.Sleep(time.Second * 10)

}
