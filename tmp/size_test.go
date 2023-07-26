package main

import (
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestOutputQuote(t *testing.T) {
	src := "12321%1123\"\"22231"
	t.Log(strconv.Quote(src))
}

type SizeA struct {
	A int
	B bool
}

type SizeB struct {
	C bool
	D int
}

type SizeC struct {
	SizeA
	SizeB
}

type SizeD struct {
	A int
	B bool
	C bool
	D int
}

func TestLogSize(t *testing.T) {
	t.Log(unsafe.Sizeof(SizeA{}))
	t.Log(unsafe.Sizeof(SizeB{}))
	t.Log(unsafe.Sizeof(SizeC{}))
	t.Log(unsafe.Sizeof(SizeD{}))
}

type InnerA struct {
	A int
	B int32
}

type InnerB struct {
	InnerA
	C int32
}

type Inner struct {
	A int
	B int32
}

type InnerD struct {
	C     int32
	inner Inner
}

//go:linkname valueInterface reflect.valueInterface
func valueInterface(reflect.Value, bool) interface{}

func TestSizeOfInner(t *testing.T) {

	{
		var a InnerD
		var v = reflect.ValueOf(&a)
		inner := v.Elem().FieldByName("inner")

		innerAddr := inner.Pointer()

		ii := valueInterface(inner, false)

		innerV := (*Inner)(unsafe.Pointer(&ii))
		innerV.A = 1
		innerV.B = 2
		t.Log(inner, a.inner, innerAddr)
	}

	{
		var a InnerD
		var v = reflect.ValueOf(&a)

		// offset 是相较于直接父结构体的偏移量, 如果多次嵌套, 需要累加
		inner, _ := v.Elem().Type().FieldByName("inner")
		var totalOffset uintptr
		for i := range inner.Index {
			totalOffset += v.Elem().Type().FieldByIndex(inner.Index[:i+1]).Offset
		}

		var innerV = (*Inner)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + totalOffset))
		innerV.A = 1
		innerV.B = 2
		t.Log(inner, a.inner)
	}

	{
		var b InnerB
		var v = reflect.ValueOf(&b)
		v.Elem().FieldByName("A").SetInt(1)
		t.Log(b)
	}

	t.Log(unsafe.Sizeof(InnerA{}))
	t.Log(unsafe.Sizeof(InnerB{}), unsafe.Offsetof(InnerB{}.C))
	t.Log(unsafe.Sizeof(InnerC{}), unsafe.Offsetof(InnerC{}.C))
	t.Log(unsafe.Sizeof(InnerD{}), unsafe.Offsetof(InnerD{}.C))
}

type InnerC struct {
	A int
	B int32
	C int32
}
