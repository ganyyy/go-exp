package main

import (
	"flag"

	"ganyyy.com/go-exp/demo/rsa/cmd"
)

var client = flag.Bool("client", false, "start client")
var addr = flag.String("addr", "", "server address")

func main() {
	flag.Parse()

	if *client || *addr != "" {
		cmd.Client(*addr)
	} else {
		cmd.Server()
	}
}
