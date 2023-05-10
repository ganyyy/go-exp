//go:build !C

package patch

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func init() {
	println("run in go mode")
}

const ()

func dlError() error {
	var reason = C.dlerror()
	if reason == nil {
		return nil
	}
	var str = C.GoString((*C.char)(reason))
	return errors.New(str)
}

func PluginOpen(path string) (unsafe.Pointer, error) {
	var cPath *C.char
	if path != "" {
		cPath = toCString(path)
		defer freeCString(cPath)
	}
	// 优化: 通过RTLD_LOCAL|RTLD_NOW, 可以将符号表的搜索范围限制在当前模块内
	// 注意: 不能使用 RTLD_DEFAULT, 替代 dlopen(NULL, ...), 因为 RTLD_DEFAULT 会导致全局符号表的搜索
	var handler = C.dlopen(cPath, C.RTLD_LOCAL|C.RTLD_NOW)
	return handler, dlError()
}

func LookupSymbol(handler unsafe.Pointer, symbolName string) (unsafe.Pointer, error) {
	var cName = toCString(symbolName)
	defer freeCString(cName)
	var pointer = C.dlsym(handler, cName)
	return pointer, dlError()
}

func PluginClose(handler unsafe.Pointer) error {
	C.dlclose(handler)
	return dlError()
}
