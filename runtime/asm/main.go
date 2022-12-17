package main

func myadd(int, int) int
func mtest()

//go:noinline
func add2(a, b, d int) int {
	c := 0x11111111 + d
	return a*b + c
}

func main() {
	var a = 100
	var b = 200
	a = myadd(a, b)
	mtest()
	println(a, b, add2(a, b, a+b))
}
