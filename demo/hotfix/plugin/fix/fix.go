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
	println("lissss 10089:", v)
	return 0
}

type MyData struct {
	A, B, C, D int
}

func (m *MyData) SetA(a int) {
	m.B = a + 10
	m.C = a + 20
	println("in plugin MyData_SetA:", a, m.B, m.C)
}
