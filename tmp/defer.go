package main

import (
	"log"
	_ "unsafe"
)

func main1() {
	var counter int
	defer println(counter)
	counter = 5
	println(counter)
}
func main2() {
	f := func() (r int) {
		defer func(r int) {
			r = r + 5
		}(r)
		return 1
	}
	println(f())
}
func main3() {
	f := func() (r int) {
		t := 5
		defer func() {
			t = t + 5
		}()
		return t
	}
	println(f())
}
func main4() {
	f := func() (r int) {
		defer func() {
			r = r + 1
		}()
		return 5
	}
	println(f())
}

func overload() (c int) {
	var a int
	{
		var a uint
		_ = a
	}
	_ = a

	return
}

//go:linkname getargp runtime.getargp
func getargp() uintptr

func main() {
	var deferRecover = func() {
		println("defer", getargp())
		if err := recover(); err != nil {
			log.Println("panic: ", err)
		}
	}
	func() {
		defer deferRecover()
		println("func1", getargp())
		var v = 0
		println(1 / v)
	}()
	func() {
		defer func() {
			println("func2 defer", getargp())
			deferRecover()
		}()
		println("func2", getargp())
		var v = 0
		println(1 / v)
	}()
}
