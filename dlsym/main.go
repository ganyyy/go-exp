package main

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

#include <stdio.h>

static void* pluginMainLookup(const char* name, char** err) {
	// NULL 表示查找当前可执行文件的符号
	void* r = dlsym((void*)NULL, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}
*/
import "C"

import (
	"plugin"
	"reflect"
	"syscall"
	"unsafe"
)

func main() {
	// 必须要先执行 dlopen, 才可以查找符号表
	var handler, _ = plugin.Open("./plugin.so")

	var sum, _ = handler.Lookup("Sum")

	println(sum.(func(int) int)(100))

	println(Add(Num, 100), ForeachAdd(Num))

	var name = []byte("main.Num")
	var cErr *C.char
	var num = C.pluginMainLookup((*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if num == nil {
		println(C.GoString(cErr))
	} else {
		var newNum = (*int)(num)
		*newNum = 200
		println(*newNum, Num)
		println(uintptr(num), uintptr(unsafe.Pointer(&Num)))
	}
	ShowFuncAddr(ForeachAdd, "main.ForeachAdd")
}

func ShowFuncAddr(fun interface{}, name string) {
	var cErr *C.char
	println("func:", name)
	var funcName = []byte(name)
	var addFunc = C.pluginMainLookup((*C.char)(unsafe.Pointer(&funcName[0])), &cErr)
	if addFunc == nil {
		println(C.GoString(cErr))
	} else {
		println("dlsym:", addFunc)
	}

	if name == "main.ForeachAdd" {
		println(ForeachAdd(10))

		// addFunc 相当于拿到了 函数的地址, 进行替换跳转
		var value = reflect.ValueOf(TotalSum)
		var toAddr = (*funcVal)(unsafe.Pointer(&value)).ptr

		var toBytes = jmpToGoFn(uintptr(toAddr))

		copyToLocation(uintptr(addFunc), toBytes)

		println(ForeachAdd(10))
	}

	// var val = reflect.ValueOf(fun)
	// var addr = (*funcVal)(unsafe.Pointer(&val)).ptr
	// println("reflect:", addr)
}

type funcVal struct {
	_   uintptr
	ptr unsafe.Pointer
}

func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}

func rawMemoryAccess(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

func jmpToGoFn(to uintptr) []byte {
	return []byte{
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}

var Num int = 30

//go:noinline
func Add(a, b int) int {
	return a + b
}

//go:noinline
func ForeachAdd(a int) int {
	var ret int
	for i := 0; i < a; i++ {
		ret += i * 5
	}
	return ret
}

//go:noinline
func TotalSum(a int) (ret int) {
	for i := 0; i < a*2; i++ {
		ret += 10
	}
	return ret
}

type Stu struct {
	Name string
}

//go:noinline
func (s *Stu) SetName(newName string) {
	s.Name = newName
}
