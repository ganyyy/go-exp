package main

import (
	"fmt"
	"reflect"
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
	Len     int
}

//go:noinline
func (p *PointerStruct2) SetValue(v []byte) {
	var head = (*reflect.SliceHeader)(unsafe.Pointer(&v))
	p.UintPtr = head.Data
	p.Len = head.Len
}

func (p *PointerStruct2) String() string {
	return fmt.Sprintf("%X", p.UintPtr)
}

//go:noinline
func NewPointStruct() *PointerStruct {
	return &PointerStruct{}
}

//go:noinline
func NewPointStruct2() *PointerStruct2 {
	return &PointerStruct2{}
}

func Pointer() {
	var i struct {
		_ int
		p *PointerStruct
	}
	var i2 PointerStruct2

	var pi = NewPointStruct()
	var pi2 = NewPointStruct2()

	i.p = pi
	i2.UintPtr = pi2.UintPtr

	_ = i
	_ = i2
}
