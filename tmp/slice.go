package main

import "unsafe"

type SliceData struct {
	Age  int
	Name string
}

var SliceAAA struct {
	S []*SliceData
}

//go:noinline
func ParallelAppend() {

	SliceAAA.S = make([]*SliceData, 0, 100)
	SliceAAA.S = append(SliceAAA.S, &SliceData{
		Age:  100,
		Name: "123",
	})
}

type NilStruct1 struct {
	Name string
}

//go:noinline
func (n *NilStruct1) GetName() string {
	if n == nil {
		return ""
	}
	return n.Name
}

type NilStruct2 struct {
	_ int
	*NilStruct1
}

func GetSliceAppend() {
	ParallelAppend()
	var s = SliceAAA.S
	var s2 = s
	var name = s[0]
	print(name.Age, len(s2), cap(s2))
}

func CallNilMethod() {
	var n = new(NilStruct2)
	n.NilStruct1 = new(NilStruct1)
	print(n.GetName())
	n.Name = "123"
}

func StackOver() {
	var a int
	println("a:%v, %p", a, (unsafe.Pointer(&a)))

	var tt [10000]int

	for i := 0; i < 10000; i++ {
		var t [100]int
		t[0] += 100
	}

	_ = tt[0]
	println("a:%v, %p", a, (unsafe.Pointer(&a)))
}
