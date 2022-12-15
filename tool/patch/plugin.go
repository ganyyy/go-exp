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

func dlError() error {
	var reason = C.dlerror()
	if reason == nil {
		return nil
	}
	var str = C.GoString((*C.char)(reason))
	return errors.New(str)
}

func PluginOpen(path string) (unsafe.Pointer, error) {
	var cPath = toCString(path)
	defer freeCString(cPath)
	var handler = C.dlopen(cPath, C.RTLD_GLOBAL|C.RTLD_NOW)
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
