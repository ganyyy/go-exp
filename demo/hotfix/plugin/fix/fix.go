package fix

import "math/rand"

//go:noinline
func Show(a, b, c, d, e, f, g int) {
	var v = rand.Int()
	println("in plugin 8978:", a, b, c, d, e, f, g, v)
}

//go:noinline
func Sum3(src []int) int {
	var v = rand.Int()
	println("lissss 999:", v)
	return 0
}
