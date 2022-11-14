package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	"ganyyy.com/go-exp/rpc/grpc/logger"
	"ganyyy.com/go-exp/rpc/grpc/proto"
)

type Server struct {
	proto.UnimplementedGreeteServer
}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("[SayHello] recv %v", req.String())
	return &proto.HelloResponse{
		Message: "Server" + req.GetName(),
	}, nil
}

func (s *Server) HelloStream(stream proto.Greete_HelloStreamServer) error {

	var md, _ = metadata.FromIncomingContext(stream.Context())
	id := md["id"]
	log.Printf("client %v connect", id)
	defer func() {
		log.Printf("client %v disconnect", id)
	}()

	var msgChan = make(chan *proto.HelloRequest, 100)

	var done = make(chan struct{})

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Printf("%v recv error %v", id, err)
				} else {
					log.Printf("%v recv close", id)
				}
				close(done)
				break
			}
			log.Printf("%v recv %v", id, msg)
			select {
			case msgChan <- msg:
			default:
				log.Printf("%v send msg full", id)
			}
		}
	}()
end:
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				break end
			}
			if msg == nil {
				continue
			}
			err := stream.Send(&proto.HelloResponse{
				Message: msg.GetName(),
			})
			if err != nil {
				log.Printf("%v send error %v", id, err)
				break end
			}
		case _, ok := <-done:
			if !ok {
				log.Printf("%v send done", id)
				break end
			}
		}
	}
	// 服务器返回就意味着流的结束
	return nil
}

var (
	port = flag.Int("port", 10086, "The server port")
)

func main() {
	logger.SetGRPCLogger()

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
		MaxConnectionIdle: time.Second * 10, // Idle连接的存活最大时长. 最近一次请求/连接建立成功中的最大值开始计算. 通过GoAway关闭
		// MaxConnectionAge:      time.Second * 15, // 任何链接存活超过30s, 发送GoAway. 这个值有 +-10% 的浮动
		MaxConnectionAgeGrace: time.Second * 5, // 超过MaxConnectionAge(已发送GoAway), 剩余的最长等待的时间
		Time:                  time.Second * 5, // 向空闲客户端发送Ping的时间间隔
		Timeout:               time.Second,     // 空闲链接回复Ping的超时
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
