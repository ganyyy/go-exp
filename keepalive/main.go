//go:build amd64 && linux

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var c = flag.Bool("c", false, "client mode")

const ADDR = ":8899"

func ErrCheck(err error, reason string) {
	if err == nil {
		return
	}
	log.Fatalf("[ERR] %s: %v", reason, err)
}

func main() {
	flag.Parse()

	var ctx, cancel = context.WithCancel(context.TODO())
	go func() {
		if *c {
			RunClient(ctx)
		} else {
			RunServer(ctx)
		}
	}()

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT)

	<-sigChan

	cancel()

	time.Sleep(time.Second)

	var mode = "Server"
	if *c {
		mode = "Client"
	}

	log.Printf("%v Exit!", mode)
}
