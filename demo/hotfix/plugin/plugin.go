package main

import (
	"ganyyy.com/go-exp/demo/hotfix/plugin/fix"

	_ "unsafe"
)

var (
	Show = fix.Show
	_    = fix.Sum3
)

func Sum3(src []int) int {
	return fix.Sum3(src)
}

var Version string

var Data MyData

func init() {
	println("load plugin:", Version)
}

//go:noinline
func Add(a, b, c int) int {
	return a * b * c * 2
}

// Exchange: newFunc, oldFunc
var Exchange = map[string]string{
	"Sum3":        "main.Min",
	"MyData_SetA": "main.(*iData).SetA",
	// "Add": "main.Empty",
}

type MyData struct {
	A, B, C, D int
}

func MyData_SetA(m *MyData, a int) {
	m.B = a + 10
	m.C = a + 20
	println("in plugin MyData_SetA:", a, m.B, m.C)
}
