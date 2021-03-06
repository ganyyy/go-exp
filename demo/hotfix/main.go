package main

import (
	"math"
	"time"

	"ganyyy.com/go-exp/demo/hotfix/update"
)

//go:noinline
func Min(src []int) int {
	var ret = math.MaxInt32
	for _, v := range src {
		if ret > v {
			ret = v
		}
	}
	return ret
}

//go:noinline
func Show(a, b, c, d, e, f, g int) {
	println("in main:", a, b, c, d, e, f, g)
}

//go:noinline
func Empty(a, b, c int) int {
	return a + b + c
}

func main() {
	go update.RunUpdateMonitor()
	var src = []int{1, 2, 3, 4, 5}
	for {
		time.Sleep(time.Second)
		println("main src: ", Min(src))
	}
}
