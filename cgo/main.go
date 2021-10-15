package main

/*
#include "./plugin.h"
*/
import "C"
import "log"

//go:generate go tool cgo main.go

////export GoF
//func GoF(arg1, arg2 int, args string) int64 {
//	log.Println(arg1, arg2, args)
//	return int64(arg1 + arg2)
//}

func main() {

	var sock = C.tcpSocket()

	log.Println(sock)
}
