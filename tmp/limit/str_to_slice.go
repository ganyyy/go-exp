package main

import (
	"reflect"
	"unsafe"
)

func ErrStringToSlice(str string) []byte {
	// 这里的 * 相当于一次copy
	var header = *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}))
}

func ErrStringToSlice2(str string) []byte {
	var header = (*reflect.StringHeader)(unsafe.Pointer(&str))
	var sliceHeader reflect.SliceHeader
	sliceHeader.Data = header.Data // 这里相当于一次copy
	sliceHeader.Len = header.Len
	sliceHeader.Cap = header.Len
	return *(*[]byte)(unsafe.Pointer(&sliceHeader)) // 这里又是一次拷贝
}

// 正确的转换

// 推荐这种
func OkStringToSlice(str string) []byte {
	var header = (*reflect.StringHeader)(unsafe.Pointer(&str))
	var sb []byte
	var sliceHeader = (*reflect.SliceHeader)(unsafe.Pointer(&sb))
	sliceHeader.Data = header.Data
	sliceHeader.Len = header.Len
	sliceHeader.Cap = header.Len
	return sb
}

func OkStringToSlice2(str string) []byte {
	var header = (*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}))
}

func Test() {
	var sb = []byte("123213213")
	var src = string(sb)
	var bs1 = ErrStringToSlice(src)
	// var bs2 = ErrStringToSlice2(src)
	// var bs3 = OkStringToSlice(src)
	var bs4 = OkStringToSlice2(src)

	println(bs1, bs4)
}
