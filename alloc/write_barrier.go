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
}

//go:noinline
func (p *PointerStruct2) SetValue(v []byte) {
	var head = (*reflect.SliceHeader)(unsafe.Pointer(&v))
	p.UintPtr = head.Data
}

func (p *PointerStruct2) String() string {
	return fmt.Sprintf("%X", p.UintPtr)
}
