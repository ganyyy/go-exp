package main

import (
	"context"
	"log"
	"net"

	"golang.org/x/sys/unix"
)

func RunServer(done context.Context) {
	var ln, err = net.Listen("tcp", ADDR)
	ErrCheck(err, "Create Listen "+ADDR)

	conn, err := ln.Accept()
	ErrCheck(err, "Accept")

	log.Printf("Client %v conn!", conn.RemoteAddr())

	var tcpConn, _ = conn.(*net.TCPConn)
	var f, _ = tcpConn.File()
	// KEEP ALIVE
	unix.SetsockoptInt(int(f.Fd()), unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1)

	// 检查超时
	unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_KEEPIDLE, 10)

	// 检查间隔
	unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_KEEPINTVL, 10)

	// 等待退出
	<-done.Done()
	log.Printf("Server Exit!")
}
