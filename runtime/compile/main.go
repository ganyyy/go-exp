package main

import (
	"math/rand"
	"os"
	"strconv"
)

// go: noinline
func Add(a, b int) int {
	var ret int
	for i := 0; i < a; i++ {
		ret += i + a*b
	}
	return ret
}

// go: noinline
func Open(f string) (*os.File, error) {
	return os.OpenFile(f, os.O_CREATE|os.O_RDWR, os.ModePerm)
}

func main() {
	var a, b int
	a = rand.Intn(10)
	b = rand.Intn(20)
	ret := Add(a, b)
	f, _ := Open("local.txt")
	f.WriteString(strconv.Itoa(ret))
	_ = f.Close()

	var arr = []int{1, 2, 5: 10, 3, 3: 20, 30, 7: 29}
	println(len(arr))
}
