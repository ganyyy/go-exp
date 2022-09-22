package main

import (
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
