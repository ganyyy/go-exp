//go:build ignore

package main

import (
	"runtime"
	"time"
)

type FFF struct {
	I [1 << 20]int
}

func init() {
	var t = newItem()

	runtime.SetFinalizer(t, func(v interface{}) {
		println("release init t")
	})

	runtime.KeepAlive(t)
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
	// var st = time.Now()
	// runtime.SetFinalizer(i, (*FFF).close) // 正确的. 成员方法本质上就是接收对象作为首个参数的方法

	// runtime.SetFinalizer(i, func(x interface{}) {
	// 	runtime.SetFinalizer(i, nil) // 错误的, 这种循环引用会导致该对象一直无法释放.
	// 	println(1)
	// })

	runtime.SetFinalizer(i, func(x interface{}) {
		runtime.SetFinalizer(x, nil) // 正确的, 此时没有捕获外部变量
		println(1)
	})
}

func main() {
	testFinalizer()
	time.Sleep(time.Second * 1)
	runtime.GC()
	time.Sleep(time.Second * 1)
	time.Sleep(time.Second * 1)

}
