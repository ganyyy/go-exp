package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

type PointerStruct struct {
	Ptr unsafe.Pointer
	Len int
}

//go:noinline
func (p *PointerStruct) SetValue(v []byte) {
	var head = (*reflect.SliceHeader)(unsafe.Pointer(&v))
	p.Ptr = unsafe.Pointer(head.Data)
	p.Len = head.Len
}

func (p *PointerStruct) String() string {
	return fmt.Sprintf("%X", p.Ptr)
}

type PointerStruct2 struct {
	UintPtr uintptr
}

//go:noinline
func (p *PointerStruct2) SetValue(v []byte) {
	var head = (*reflect.SliceHeader)(unsafe.Pointer(&v))
	p.UintPtr = head.Data
}

func (p *PointerStruct2) String() string {
	return fmt.Sprintf("%X", p.UintPtr)
}

func writeBarrier() {

	var p = new(PointerStruct)
	var p2 = new(PointerStruct2)

	showTypeInfo(p)
	showTypeInfo(p2)

	var tmp = make([]byte, 1<<20)
	var tmp2 = make([]byte, 1<<20)
	tmp[1] = 'A'
	p.SetValue(tmp)
	p2.SetValue(tmp2)

	println("p:", p.String(), "\t p2:", p2.String())
	runtime.GC() // 标记
	time.Sleep(time.Second)
	runtime.GC() // 清理

	println("p:", p.String(), "\t p2:", p2.String())
	p.SetValue(make([]byte, 1<<20)) // 此时复用的是tmp2所指向的那片内存
	println("p:", p.String(), "\t p2:", p2.String())
}
