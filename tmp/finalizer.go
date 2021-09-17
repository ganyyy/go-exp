//go:build ignore

package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var i = make([]int, 1 << 20)

	i[0] = 100
	//var st = time.Now()
	runtime.SetFinalizer(&i, func(v interface{}) {
		fmt.Println(1)
	})

	time.Sleep(time.Second * 10)
}
