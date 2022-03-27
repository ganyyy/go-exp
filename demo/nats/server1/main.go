package main

import (
	"log"
	"os"
	"os/signal"

	"ganyyy.com/go-exp/demo/nats/method"
	"ganyyy.com/go-exp/demo/nats/nc"
)

const (
	ADDR = "localhost"
	PORT = 4222
)

func main() {
	var config nc.NatsClientConfig
	config.Addr = ADDR
	config.Port = PORT
	if err := nc.Init(config); err != nil {
		panic(err)
	}
	defer nc.Stop()

	var service, err = nc.NewNatsServiceModule("1", method.TestReq, method.TestPush)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case msg := <-service.Message():
			switch msg.Method {
			case method.PushTime:
				log.Printf("[INFO] Req:%v", string(msg.Data))
				var param method.PushTimeParam
				_ = nc.Decode(msg.Data, &param)
				var req method.MethReqTime
				req.Time = param.Time

				var rsp method.MethodRspTime
				service.Request(method.ReqTime, req, &rsp)
				log.Printf("PushTime:%+v, Rsp:%+v", param, rsp)
			case method.PushValue:
			default:
				log.Printf("[ERR] receive msg:%+v, %+v", msg.Method, string(msg.Data))
			}
		case <-quit:
			log.Printf("[INF] Done!")
			return
		}
	}
}
