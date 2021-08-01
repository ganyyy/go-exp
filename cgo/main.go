package main

/*
#include "./plugin.h"
*/
import "C"
import "log"

//go:generate go tool cgo main.go

func main() {

	var sock = C.tcpSocket()

	log.Println(sock)
}
