package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
)

//go:generate protoc --proto_path=. --go_out=. *.proto

func main() {
	var req Req
	req.Age = 100

	var bs, _ = proto.Marshal(&req)
	log.Println(bs)

	fmt.Printf("%[2]v, %[1]v\n", bs, string(bs))

	testBlockChan()

}

func testBlockChan() {
	var ch = make(chan int, 1)

	go func() {
		select {
		case ch <- 1:
		default:
			log.Println("send to default")
		}
	}()

	log.Println(<-ch)

}