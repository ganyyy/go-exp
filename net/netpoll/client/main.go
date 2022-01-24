package main

import (
	"go-exp/net/msg"
	net "go-exp/net/netpoll"
	"log"
	"time"

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
