package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"ganyyy.com/go-exp/rpc/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 10086, "The server port")
)

func main() {
	flag.Parse()

	var conn, err = grpc.Dial(fmt.Sprintf(":%v", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("dial error:%v", err)
	}
	defer conn.Close()

	var client = proto.NewGreeteClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rsp, err := client.SayHello(ctx, &proto.HelloRequest{
		Name: "123131",
	})
	if err != nil {
		log.Panicf("dial error:%v", err)
	}
	log.Printf("resp:%v", rsp.GetMessge())
}
