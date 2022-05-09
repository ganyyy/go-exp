package main

//go:generate go build -buildmode=plugin sum.go

//go:noinline
func Add(a, b int) int {
	println("in sum.so add")
	return a + b
}

//go:noinline
func Sub(a, b int) int {
	println("in sum.so sub")
	return a - b
}
