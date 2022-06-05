package main

import (
	"flag"

	"ganyyy.com/go-exp/demo/kcp-go/cmd"
)

var client = flag.Bool("client", false, "use client mode")

func main() {
	flag.Parse()

	if *client {
		cmd.Client()
	} else {
		cmd.Server()
	}
}
