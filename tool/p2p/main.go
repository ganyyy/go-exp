package main

import (
	"p2p/cmd"
	_ "p2p/cmd/client"
	_ "p2p/cmd/server"
)

func main() {
	cmd.Run()
}
