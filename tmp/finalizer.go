//go:build ignore

package main

import (
	"runtime"
	"time"
)

type FFF struct {
	I [1 << 20]int
}

func (f *FFF) close() {
	runtime.SetFinalizer(f, nil)
	println(1)
}

func newItem() *FFF {
	return &FFF{}
}

func testFinalizer() {
	var i = newItem()
	//var st = time.Now()
	runtime.SetFinalizer(i, (*FFF).close)
	// runtime.SetFinalizer(i, func(x interface{}) {
	// 	runtime.SetFinalizer(i, nil) // 错误的
	// 	println(1)
	// })
}

func main() {
	testFinalizer()
	time.Sleep(time.Second * 1)
	runtime.GC()
	time.Sleep(time.Second * 1)
	time.Sleep(time.Second * 1)

}
