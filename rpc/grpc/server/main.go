package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"ganyyy.com/go-exp/rpc/grpc/proto"
)

type Server struct {
	proto.UnimplementedGreeteServer
}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Messge: "Server" + req.GetName(),
	}, nil
}

var (
	port = flag.Int("port", 10086, "The server port")
)

func main() {
	flag.Parse()
	var lis, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panicf("Listen error:%v", err)
	}

	var keepPolicy = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // 客户端两次Ping之间的间隔
		PermitWithoutStream: true,            // 是否支持在没有Stream的情况下接受客户端的ping
	}

	var keepServer = keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * 10, // Idle连接的存活最大时长. 最近一次请求/连接建立成功中的最大值开始计算. 通过GoAway关闭
		MaxConnectionAge:      time.Second * 15, // 任何链接存活超过30s, 发送GoAway. 这个值有 +-10% 的浮动
		MaxConnectionAgeGrace: time.Second * 5,  // 超过MaxConnectionAge(已发送GoAway), 剩余的最长等待的时间
		Time:                  time.Second * 5,  // 向空闲客户端发送Ping的时间间隔
		Timeout:               time.Second,      // 空闲链接回复Ping的超时
	}

	var server = grpc.NewServer(
		grpc.KeepaliveParams(keepServer),
		grpc.KeepaliveEnforcementPolicy(keepPolicy),
	)
	proto.RegisterGreeteServer(server, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Panicf("failed to server %v", err)
	}
}
