//go:build !C

package patch

/*
#cgo linux LDFLAGS: -ldl -g -O2
#include <dlfcn.h>
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
	var handler = C.dlopen(toCString(path), C.RTLD_GLOBAL|C.RTLD_NOW)
	return handler, dlError()
}

func LookupSymbol(handler unsafe.Pointer, symbolName string) (unsafe.Pointer, error) {
	var pointer = C.dlsym(handler, toCString(symbolName))
	return pointer, dlError()
}

func PluginClose(handler unsafe.Pointer) error {
	C.dlclose(handler)
	return dlError()
}
