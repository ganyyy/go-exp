package main

import "math/rand"

var Version string

func init() {
	println("load plugin:", Version)
}

//go:noinline
func Sum3(src []int) int {
	var v = rand.Int()
	println("lissss 666:", v)
	return 0
}

//go:noinline
func Show(a, b, c, d, e, f, g int) {
	var v = rand.Int()
	println("in plugin 3333:", a, b, c, d, e, f, g, v)
}

//go:noinline
func Add(a, b, c int) int {
	return a * b * c * 2
}

//Exchange: newFunc, oldFunc
var Exchange = map[string]string{
	"Sum3": "main.Min",
	// "Add": "main.Empty",
}
