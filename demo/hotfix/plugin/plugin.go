package main

import (
	"reflect"
	"runtime"

	"ganyyy.com/go-exp/demo/hotfix/plugin/fix"

	_ "unsafe"
)

var Version string

func init() {
}

// Exchange: newFunc, oldFunc
var Exchange = map[string]string{
	"Sum3":        "main.Min",
	"MyData_SetA": "ganyyy.com/go-exp/demo/hotfix/common.(*Data).SetA",
	"MyFixAdd":    "main.GenAdd.func1",
	// "Add": "main.Empty",
}

// 替换函数实现

func Sum3(src []int) int {
	return fix.Sum3(src)
}

func MyData_SetA(m *fix.MyData, a int) {
	m.SetA(a)
}

var MyFixAdd = runtime.FuncForPC(reflect.ValueOf(GenAdd2Fix(10)).Pointer()).Entry()

func GenAdd2Fix(a int) func(int) int {
	return func(i int) int {
		return a
	}
}
