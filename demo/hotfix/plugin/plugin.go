package main

import (
	"ganyyy.com/go-exp/demo/hotfix/plugin/fix"

	_ "unsafe"
)

var Version string

func init() {
	println("load plugin:", Version)
}

// Exchange: newFunc, oldFunc
var Exchange = map[string]string{
	"Sum3":        "main.Min",
	"MyData_SetA": "ganyyy.com/go-exp/demo/hotfix/common.(*Data).SetA",
	// "Add": "main.Empty",
}

// 替换函数实现

func Sum3(src []int) int {
	return fix.Sum3(src)
}

func MyData_SetA(m *fix.MyData, a int) {
	m.SetA(a)
}
