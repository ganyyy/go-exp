package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"ganyyy.com/go-exp/rpc/grpc/proto"
	"google.golang.org/grpc"
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
	var server = grpc.NewServer()
	proto.RegisterGreeteServer(server, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Panicf("failed to server %v", err)
	}
}
