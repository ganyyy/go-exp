package pkg

import (
	"reflect"
	"runtime"
)

//go:noinline
func Add(a, b int) int {
	return a + b
}

func funcInfo(f interface{}) *runtime.Func {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer())
}

//go:noinline
func FuncName(f interface{}) string {
	return funcInfo(f).Name()
}

//go:noinline
func FuncEntry(f interface{}) uintptr {
	return funcInfo(f).Entry()
}
