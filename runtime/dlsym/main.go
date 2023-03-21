package main

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

static uintptr_t pluginSelf(const char* path,char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_LOCAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static uintptr_t pluginOpen(const char* path, char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static void* pluginLookup(uintptr_t h, const char* name, char** err) {
	void* r = dlsym((void*)h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}

static void* pluginMainLookup(const char* name, char** err) {
	return pluginLookup((uintptr_t)0, name, err);
}

static int pluginClose(uintptr_t h, char** err) {
	int ret = dlclose((void*)h);
	if (ret != 0) {
		*err = (char*)dlerror();
	}
	return ret;
}
*/
import "C"

import (
	"flag"
	"fmt"
	"plugin"
	"reflect"
	"runtime"
	"unsafe"

	"ganyyy.com/go-exp/patch"
	"ganyyy.com/go-exp/runtime/dlsym/pkg"
)

var (
	fName = flag.String("f", "", "func name")
)

func main() {
	flag.Parse()

	var hanler2, _ = plugin.Open("./plugin/plugin.so")

	pAddr, _ := hanler2.Lookup("FuncAddr")
	addr, _ := pAddr.(func() map[string]uintptr)

	addrMap := addr()

	f1, _ := hanler2.Lookup("Add")
	f2, _ := hanler2.Lookup("Add2")

	ff1, _ := f1.(func(int, int) int)
	ff2, _ := f2.(func(int, int) int)

	println(pkg.FuncName(TotalSum))

	println(f1, f2)
	println(ff1(10, 20), ff2(10, 20))

	fmt.Println(runtime.FuncForPC(addrMap[pkg.FuncName(pkg.Add)]).Name())

	fmt.Println(addrMap)

	patch.Patch(unsafe.Pointer(addrMap["Add"]), unsafe.Pointer(addrMap["Add2"]))
	println(uintptr(patch.FuncAddr(f1)), uintptr(patch.FuncAddr(f2)))
	println(uintptr(addrMap["Add"]), uintptr(addrMap["Add2"]))
	println(ff1(10, 20), ff2(10, 20))

}

func oldDlsym() {
	// 必须要先执行 dlopen, 才可以查找符号表
	var handler2, err = plugin.Open("./plugin.so")
	println(err)

	var sum, _ = handler2.Lookup("Sum")
	var ff = sum.((func(int) int))
	println(ff(100))

	println(Add(Num, 100), ForeachAdd(Num))

	var cErr *C.char
	// {
	// 	var name = []byte("./plugin.so")
	// 	var handle = C.pluginOpen((*C.char)(unsafe.Pointer(&name[0])), &cErr)
	// 	if handle == 0 {
	// 		panic(C.GoString(cErr))
	// 	}
	// 	var sumName = []byte("Sum")
	// 	// var val interface{}
	// 	// valp := (*[2]unsafe.Pointer)(unsafe.Pointer(&val))
	// 	var sum = C.pluginLookup(handle, (*C.char)(unsafe.Pointer(&sumName[0])), &cErr)
	// 	// valp[1] = unsafe.Pointer(&sum)
	// 	// var ret = val.(func(int) int)(100)
	// 	println(sum, C.GoString(cErr))
	// 	var ret = C.pluginClose(handle, &cErr)
	// 	println(ret, C.GoString(cErr))
	// 	return
	// }

	var handler = C.pluginSelf((*C.char)(unsafe.Pointer(uintptr(0))), &cErr)
	if handler == 0 {
		panic("cannot open plugin")
	}

	var name = append([]byte("main.Num"), 0)
	var num = C.pluginLookup(handler, (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if num == nil {
		println(C.GoString(cErr))
	} else {
		var newNum = (*int)(num)
		*newNum = 200
		println(*newNum, Num)
		println(uintptr(num), uintptr(unsafe.Pointer(&Num)))
	}
	fname := "main.ForeachAdd"
	if len(*fName) != 0 {
		fname = *fName
	}
	ShowFuncAddr(uint64(handler), nil, fname)
}

func ShowFuncAddr(handle uint64, fun interface{}, name string) {
	var cErr *C.char
	println("func:", name)
	var funcName = append([]byte(name), 0)
	var addFunc = C.pluginLookup(C.ulong(handle), (*C.char)(unsafe.Pointer(&funcName[0])), &cErr)
	if addFunc == nil {
		println(C.GoString(cErr))
		return
	} else {
		println("dlsym:", addFunc)
	}

	if name == "main.ForeachAdd" {
		println(ForeachAdd(10))

		// addFunc 相当于拿到了 函数的地址, 进行替换跳转
		var value = reflect.ValueOf(fun)
		var toAddr = (*funcVal)(unsafe.Pointer(&value)).ptr

		var toBytes = jmpToGoFn(uintptr(toAddr))

		copyToLocation(uintptr(addFunc), toBytes)

		println(ForeachAdd(10))
	}

	// var val = reflect.ValueOf(fun)
	// var addr = (*funcVal)(unsafe.Pointer(&val)).ptr
	// println("reflect:", addr)
}

var Num int = 30

//go:noinline
func Add(a, b int) int {
	return a + b
}

//go:noinline
//export ForeachAdd
func ForeachAdd(a int) int {
	var ret int
	for i := 0; i < a; i++ {
		ret += i * 5
	}
	return ret
}

//go:noinline
//export TotalSum
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
