//go:build amd64 && linux

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
	_ = unix.SetsockoptInt(int(f.Fd()), unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1)

	// 检查超时
	_ = unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_KEEPIDLE, 10)

	// 检查间隔
	_ = unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_KEEPINTVL, 5)

	// 检查次数
	_ = unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_KEEPCNT, 1)

	// 设置USER_TIMEOUT: TCP_KEEPIDLE + TCP_KEEPINTVL * TCP_KEEPCNT
	_ = unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, 15*1000)

	// 等待退出
	<-done.Done()
	log.Printf("Server Exit!")
}
