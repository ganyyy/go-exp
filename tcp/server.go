// +build ignore

package main

import (
	"go-exp/helper"
	"log"
	"net"
)

func main() {

	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", "0.0.0.0:9999")
	helper.PanicIfErr("listen", err)

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Printf("conn accept error:%v", err)
			continue
		}

		log.Printf("client addr:%v", conn.RemoteAddr())
		conn.Close()
	}

}
