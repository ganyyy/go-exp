package main

import (
	"context"
	"go-exp/net/msg"
	net "go-exp/net/netpoll"
	"log"

	"github.com/cloudwego/netpoll"
)

type rpcServer struct{}

var key = struct{}{}

func (s *rpcServer) Run(network string, address string) error {
	listener, err := netpoll.CreateListener(network, address)
	if err != nil {
		return err
	}
	var idx int
	eventLoop, err := netpoll.NewEventLoop(s.handler,
		netpoll.WithOnPrepare(func(connection netpoll.Connection) context.Context {
			idx++
			log.Printf("%v connect!", connection.RemoteAddr())
			return context.WithValue(context.TODO(), key, idx)
		}),
	)
	if err != nil {
		return err
	}
	return eventLoop.Serve(listener)
}

func (s *rpcServer) handler(ctx context.Context, conn netpoll.Connection) (err error) {
	reader, writer := conn.Reader(), conn.Writer()

	var req msg.Message

	err = net.Decode(reader, &req)
	if err != nil {
		return err
	}
	var idx = ctx.Value(key)
	log.Printf("%v send value %v", idx, req.Message)
	req.Message += "[Rsp]"
	return net.Encode(writer, &req)
}

func main() {
	log.Panicln((&rpcServer{}).Run("tcp", "localhost:8899"))
}
