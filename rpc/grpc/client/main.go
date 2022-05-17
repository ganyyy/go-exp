package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"ganyyy.com/go-exp/rpc/grpc/proto"
)

var (
	port = flag.Int("port", 10086, "The server port")
)

func main() {
	flag.Parse()

	var keepParam = keepalive.ClientParameters{
		Time:                10 * time.Second, // 没有活跃的情况下, 最长多久发一次心跳包
		Timeout:             time.Second,      // 心跳包超过多久会认为链接已断开
		PermitWithoutStream: true,             // 没有激活的流, 是否允许发送心跳包?(可以理解为, 是针对RPC的链接保活)
	}

	var conn, err = grpc.Dial(
		fmt.Sprintf(":%v", *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepParam),
	)
	if err != nil {
		log.Panicf("dial error:%v", err)
	}
	defer conn.Close()

	var client = proto.NewGreeteClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var doCall = func() {
		rsp, err := client.SayHello(ctx, &proto.HelloRequest{
			Name: "123131",
		})
		if err != nil {
			log.Printf("dial error:%v", err)
		} else {
			log.Printf("resp:%v", rsp.GetMessge())
		}
	}

	doCall()
	time.Sleep(time.Second * 15)
	doCall()
}
