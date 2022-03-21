package main

import (
	"log"
	"runtime/debug"
)

//go:generate go build -ldflags "-X 'main.Version=$(date)'" -o info

var Version = "123"

func main() {

	var info, _ = debug.ReadBuildInfo()
	log.Printf("%+v", info)

	log.Println("Version:", Version)
}
