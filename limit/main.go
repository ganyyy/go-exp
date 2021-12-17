package main

import (
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

func main() {
	TestGC()
}

type uintptr3 [3]uintptr
type uintptr2 [2]uintptr

func TestGC() {
	var src = make([]byte, 1024*1024)

	src[100] = 'c'

	runtime.SetFinalizer(&src, func(_ interface{}) {
		println("release src")
	})

	var str = string(src)
	runtime.SetFinalizer(&str, func(_ interface{}) {
		println("release str")
	})
	println((*uintptr3)(unsafe.Pointer(&src))[0], (*uintptr2)(unsafe.Pointer(&str))[0])
	go func() {
		var ret = UnsafeToString(str)
		ret[199] = '1'
		println("ret:", (*uintptr3)(unsafe.Pointer(&ret))[0])
		var tmp = make([]byte, 1024*1024)
		var tmp2 = make([]byte, 1024*1024)
		println("tmp:", (*uintptr3)(unsafe.Pointer(&tmp))[0], "tmp2:", (*uintptr3)(unsafe.Pointer(&tmp2))[0])
		tmp2[198] = '2'
		tmp2[199] = 'A'
		println(tmp2[198], ret[199])
	}()
	time.AfterFunc(time.Second, func() {
		runtime.GC()
	})
	time.Sleep(time.Second * 6)
}

func UnsafeToString(str string) []byte {
	// return nil
	var head = *(*reflect.StringHeader)(unsafe.Pointer(&str)) // 这样写是不对的. 因为发生了一次 uintptr 的拷贝
	// var head = (*reflect.StringHeader)(unsafe.Pointer(&str))

	// 这种写法也是不对的. 因为同样发生了 uintptr 的拷贝. reflect.SliceHeader.Data 不是一个合法的指针
	// var sh reflect.SliceHeader
	// sh.Data = head.Data
	// sh.Cap = head.Len
	// sh.Len = head.Len
	time.Sleep(time.Second * 3)
	runtime.GC()
	// return *(*[]byte)(unsafe.Pointer(&sh))

	// 这样写是正确的. 因为先存在了一个合法的 sliceStruct.data, 而后进行的赋值
	// return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
	// 	Data: head.Data,
	// 	Len:  head.Len,
	// 	Cap:  head.Len,
	// }))

	var sb []byte
	var sh = (*reflect.SliceHeader)(unsafe.Pointer(&sb))
	sh.Data = head.Data
	sh.Cap = head.Len
	sh.Len = head.Len

	return sb
}
