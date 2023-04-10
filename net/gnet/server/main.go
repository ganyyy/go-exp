package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/panjf2000/gnet"
)

type echoServer struct {
	gnet.EventServer
}

// AfterWrite implements gnet.EventHandler
func (*echoServer) AfterWrite(c gnet.Conn, b []byte) {
	log.Printf("[%s] write %s", c.RemoteAddr(), b)
}

// OnInitComplete implements gnet.EventHandler
func (*echoServer) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Printf("server init success")
	return gnet.None
}

// OnOpened implements gnet.EventHandler
func (*echoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("[%s] connect!", c.RemoteAddr())
	return []byte("hello!\n"), gnet.None
}

// PreWrite implements gnet.EventHandler
func (*echoServer) PreWrite(c gnet.Conn) {
	log.Printf("[%s] pre write data", c.RemoteAddr())
}

func (e *echoServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("[%s] is close. err %v", c.RemoteAddr(), err)
	return gnet.None
}

func (e *echoServer) React(packet []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	return append([]byte("recv:"), packet...), gnet.None
}

func (e *echoServer) OnShutdown(svr gnet.Server) {
	log.Printf("all connection:%v", svr.CountConnections())
}

func main() {
	const addr = "udp://127.0.0.1:9999"
	go func() {
		err := gnet.Serve(&echoServer{}, addr,
			gnet.WithLockOSThread(true),
			gnet.WithMulticore(true),
			gnet.WithReusePort(true),
		)
		log.Fatalf("serve error %v", err)
	}()

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	log.Printf("recv %v", sig)
	log.Printf("stop %v error %v", addr, gnet.Stop(context.Background(), addr))
}
