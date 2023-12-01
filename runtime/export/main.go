package main

import (
	"fmt"
	_ "net/http/pprof"
	"reflect"
	"unsafe"

	"ganyyy.com/go-exp/runtime/export/package_b"
)

type Outer struct {
	Name string
	Age  int
	Tmp  [1 << 20]int
}

func TestExchangeGC() *Outer {
	var inner = package_b.GetInner("1234", 1234)
	var outer = *(**Outer)(unsafe.Pointer(&inner))
	return outer
}

//go:linkname Add ganyyy.com/go-exp/runtime/export/package_b.Add
func Add(a, b int) int

type mainStu struct {
	name string
}

//go:linkname setName ganyyy.com/go-exp/runtime/export/package_b.(*stu).setName
func setName(stu *mainStu, newName string)

func main() {
	fmt.Println(Add(10, 20))

	package_b.ShowFuncAddr(Add)
	package_b.ShowFuncAddr(package_b.Add)

	var funcAdd1 = reflect.ValueOf(Add)
	var funcAdd2 = reflect.ValueOf(package_b.Add)

	var rf1 = (*package_b.FuncVal)(unsafe.Pointer(&funcAdd1))
	var rf2 = (*package_b.FuncVal)(unsafe.Pointer(&funcAdd2))

	fmt.Println(rf1, rf2)

	var stu = package_b.NewStu()
	fmt.Println(stu)
	setName(*(**mainStu)(unsafe.Pointer(&stu)), "2344")
	fmt.Println(stu)
}
