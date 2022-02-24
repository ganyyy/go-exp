//go:build amd64 && linux

package main

import (
	"context"
	"log"
	"net"

	"golang.org/x/sys/unix"
)

func RunClient(done context.Context) {
	var conn, err = net.Dial("tcp", ADDR)
	ErrCheck(err, "Client Dial")

	var f, _ = conn.(*net.TCPConn).File()

	// 关闭?
	unix.SetsockoptInt(int(f.Fd()), unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0)

	_ = conn
	<-done.Done()
	log.Printf("Client Exit!")
}
