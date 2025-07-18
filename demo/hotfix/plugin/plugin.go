package main

import (
	"fmt"
	"reflect"
	"runtime"

	"ganyyy.com/go-exp/demo/hotfix/plugin/fix"
	"ganyyy.com/go-exp/demo/hotfix/plugin/locate"

	_ "unsafe"
)

var Version string

func init() {
	fmt.Println("Plugin version:", Version)
	loc, err := locate.Locate()
	fmt.Println("Plugin self path:", loc, "Error:", err)
	if err == nil {
		err = locate.DecFunc(loc)
		fmt.Println("Decryption result:", err)
	}
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
		fix.GlobalData++
		// fmt.Println("in GenAdd2Fix:", a, i, fix.GlobalData)
		return a
	}
}
