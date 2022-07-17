package main

import (
	"flag"

	"ganyyy.com/go-exp/demo/rsa/cmd"
)

var client = flag.Bool("client", false, "start client")

func main() {
	flag.Parse()

	if *client {
		cmd.Client()
	} else {
		cmd.Server()
	}
}
