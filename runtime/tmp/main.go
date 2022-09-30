package main

import "ganyyy.com/go-exp/runtime/tmp/internal"

var (
	data int
)

func init() {
	data = Add(100, 200)
}

func main() {
	println(data)
	println(Add(100, 200))
	println(internal.Add(100, 200))
	println(add(100, 200))
	println(internal.Data)
}

// func main() {
// 	p, err := plugin.Open("./plugin/plugin.so")
// 	if err != nil {
// 		panic(err)
// 	}
// 	add, err := p.Lookup("Add")
// 	if err != nil {
// 		panic(err)
// 	}
// 	println(add.(func(int, int) int)(100, 200))
// }

//go:noinline
func Add(a, b int) int {
	return a + b
}

//go:noinline
func add(a, b int) int {
	return a + b
}

//go:noinline
func sub(a, b int) int {
	return a - b
}
