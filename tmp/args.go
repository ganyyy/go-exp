//go:build ignore
// +build ignore

package main

import (
	"fmt"
)

func printArgs(a int) (a2 int) {
	return a
}

func emptyStruct() {
	var a, b struct{}
	fmt.Printf("%p, %p\n", &a, &b)
}

func main() {
	var a struct{}
	var b struct{}
	var d = 100
	var c struct{}
	fmt.Printf("%p\n", &a)
	println(&b, &c, d)
	emptyStruct()
}
