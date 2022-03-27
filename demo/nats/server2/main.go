package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

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

	var service, err = nc.NewNatsServiceModule("1", method.TestPush, method.TestReq)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var tick = time.NewTicker(time.Second)

	for {
		select {
		case msg := <-service.Message():
			switch msg.Method {
			case method.ReqTime:
				var req method.MethReqTime
				_ = nc.Decode(msg.Data, &req)
				var rsp method.MethodRspTime
				rsp.Time = strconv.Itoa(int(req.Time))
				msg.Response(rsp)
				log.Printf("[INF] Receive Request:%+v, Response:%+v", req, rsp)
			case method.ReqVal:
			}
		case <-tick.C:
			var push method.PushTimeParam
			push.Time = time.Now().Unix() + rand.Int63n(100)
			log.Printf("[INF] Start Push %+v", push)
			service.Push(method.PushTime, push)
		case <-quit:
			log.Printf("[INF] Done!")
			return
		}
	}
}
