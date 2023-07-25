package main

import (
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"

	"ganyyy.com/go-exp/demo/hotfix/common"
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
	go func() { _ = http.ListenAndServe("localhost:8899", nil) }()
	var src = []int{1, 2, 3, 4, 5}
	var data common.Data
	var idx int
	for {
		time.Sleep(time.Second)
		data.SetA(idx)
		idx++
		fmt.Printf("main src:  %v, %+v\n", Min(src), data)
	}
}
