//go:build ignore

package main

import (
	"runtime"
	"time"
	"unsafe"
)

type GCStruct struct {
	I [1 << 20]int
}

func NewStruct() *GCStruct {
	return &GCStruct{}
}

//go:noinline
func byteToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func stackAlloc() {
	var a = []byte{1, 2, 3}
	var b = byteToString(a)
	var d = a
	var c = *(*string)(unsafe.Pointer(&a))
	println(unsafe.Pointer(&a), unsafe.Pointer(&b), unsafe.Pointer(&c), unsafe.Pointer(&d))
}

func main() {
	noKeepAlive()
	hasKeepAlive()
	hasKeepAlive2()
}

func noKeepAlive() {
	var i = NewStruct()
	runtime.SetFinalizer(i, func(v interface{}) {
		println("Finalizer noKeepAlive")
	})
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
}

func hasKeepAlive() {
	var i = NewStruct()

	runtime.SetFinalizer(i, func(v interface{}) {
		println("Finalizer hasKeepAlive")
	})

	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
	runtime.KeepAlive(i)
}

func hasKeepAlive2() {
	var i = NewStruct()

	runtime.KeepAlive(i)
	runtime.SetFinalizer(i, func(v interface{}) {
		println("Finalizer hasKeepAlive2")
	})
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
	runtime.KeepAlive(i)
}
