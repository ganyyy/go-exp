//go:build ignore

package main

import (
	"fmt"
	"reflect"
	"runtime"
)

type Stu struct {
	Name string
	Age  int
}

type SAA struct {
	A, B, C byte
}

type StuPtr struct {
	SayHello
	Ptr  *byte
	Val3 [1 << 16]int64
	A    int
}

func (s *Stu) String() string {
	return ""
}

type SayHello interface {
	Hello() string
}

func main() {
	var a Stu
	var b = &a

	var aa StuPtr
	var ba = StuPtr{}

	var a = new(int)
	runtime.SetFinalizer(a, func(val interface{}) {

	})

	//ba.Hello()

	println(ba == aa)

	var aat = reflect.TypeOf(ba)
	println(aat.String())

	var at = reflect.TypeOf(a)
	var bt = reflect.TypeOf(b)

	println(at.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()))
	println(bt.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()))
	println(at.String())
	println(bt.String())

	println(reflect.TypeOf(struct {
		name string
	}{}).PkgPath())
	println(at.PkgPath())

	var s SAA
	var sst = reflect.TypeOf(s)
	println(sst.Align(), sst.Size())
}
