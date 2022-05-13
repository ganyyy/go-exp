//go:build ignore
// +build ignore

package main

import (
	"log"
	"syscall"

	"ganyyy.com/go-exp/helper"
)

func main() {
	var srvFD, cliFD int
	var err error

	srvFD, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	helper.PanicIfErr("create socket", err)

	syscall.SetsockoptInt(srvFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	syscall.Bind(srvFD, &syscall.SockaddrInet4{
		Port: 8899,
		Addr: [4]byte{
			127, 0, 0, 1,
		},
	})

	syscall.Listen(srvFD, 5)

	var addr syscall.Sockaddr
	cliFD, addr, _ = syscall.Accept(srvFD)
	log.Println(addr.(*syscall.SockaddrInet4))

	syscall.SetNonblock(cliFD, true)

	syscall.SetsockoptTimeval(cliFD, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{
		Sec: 1,
	})

	var buf [100]byte
	for i := 0; i < 10; i++ {
		n, _ := syscall.Read(cliFD, buf[:])
		if n < 0 {
			continue
		}
		log.Println(string(buf[:n]))
	}

}
