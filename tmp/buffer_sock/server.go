// +build ignore

package main

import (
	"log"
	"net"
	"time"
)

func main() {
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", ":9988")
	if err != nil {
		log.Panicf("dial error:%v", err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 3)
		conn.Write([]byte("hello world!"))
	}
	time.Sleep(time.Second * 2)
	err = conn.Close()
	if err != nil {
		log.Panicf("close conn error:%v", err)
	}
}
