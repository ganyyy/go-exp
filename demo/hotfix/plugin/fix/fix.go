package fix

import (
	"math/rand"
	_ "unsafe"

	_ "ganyyy.com/go-exp/demo/hotfix/common"
)

//go:noinline
func Show(a, b, c, d, e, f, g int) {
	var v = rand.Int()
	println("in plugin 8978:", a, b, c, d, e, f, g, v)
}

//go:linkname GlobalData main.globalData
var GlobalData int

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
	SetC(m, a+30)
	println("in plugin MyData_SetA:", a, m.B, m.C)
}

//go:linkname SetC ganyyy.com/go-exp/demo/hotfix/common.(*Data).setC
func SetC(*MyData, int)

/*
// 这是一个错误的示例, 因为 localname 对应的是一个方法, 而不是一个函数
//go:linkname (*MyData).SetC ganyyy.com/go-exp/demo/hotfix/common.(*Data).setC
func (*MyData) SetC(int)

*/
