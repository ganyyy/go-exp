package main

/*
#include "./plugin.h"
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

//go:generate go tool cgo main.go

func main() {
	//call c function
	C.helloFromC()

	var ret = C.addFromC(C.int(100), C.int(200))

	log.Println(ret)

	var cs = C.CString("12323412")
	C.print_str(cs)
	C.free(unsafe.Pointer(cs))
}

//export HelloFromGo
func HelloFromGo() {
	fmt.Printf("Hello from Go!\n")
}

//export Add
func Add(a, b int) int {
	return a + b
}

//export HelloByGo
func HelloByGo(name string) *C.char {
	return C.CString("greeting " + name)
}
