package main

import (
	"log"
	"time"

	net "ganyyy.com/go-exp/net/netpoll"

	"ganyyy.com/go-exp/net/msg"

	"github.com/cloudwego/netpoll"
)

func main() {
	var conn, _ = netpoll.DialConnection("tcp", "localhost:8899", time.Second)

	var reader, writer = conn.Reader(), conn.Writer()

	var m msg.Message

	m.Message = "1232131231"

	net.Encode(writer, &m)

	net.Decode(reader, &m)

	log.Printf("value:%v", m.Message)
}
