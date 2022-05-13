// +build ignore

package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func main() {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", ":9988")
	if err != nil {
		log.Panicf("listen error:%v", err)
	}
	var conn net.Conn
	conn, err = listener.Accept()
	if err != nil {
		log.Panicf("accept error:%v", err)
	}
	var reader = bufio.NewReaderSize(conn, 1024)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		var buf [100]byte
		var n int
		n, err = reader.Read(buf[:])
		if err != nil {
			if n != 0 {
				log.Printf("read size:%v", n)
			}
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("read timeout")
				continue
			} else {
				log.Printf("error:%v", err)
				break
			}
		}
		log.Printf("%v", string(buf[:n]))
	}

}
