//go:build ignore

package main

var base = 100

//go:noinline
func loopAdd(a int) (ret int) {
	for i := 0; i < a; i++ {
		ret += i
	}
	return ret + 100
}

func main() {
	println(loopAdd(100))
}
