//go:build C

package patch

/*
#cgo linux LDFLAGS: -ldl
#include "patch.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func init() {
	println("run in c mode")
}

func PluginOpen(path string) (unsafe.Pointer, error) {
	var cErr *C.char
	var err error
	var cPath = toCString(path)
	var handler = C.pluginOpen(cPath, &cErr)
	if handler == nil {
		err = errors.New(C.GoString(cErr))
	}
	return handler, err
}

func LookupSymbol(handler unsafe.Pointer, symbolName string) (unsafe.Pointer, error) {
	var cErr *C.char
	var err error
	var symbol = C.pluginLookup(handler, toCString(symbolName), &cErr)
	if symbol == nil {
		err = errors.New(C.GoString(cErr))
	}
	return symbol, err
}

func PluginClose(handler unsafe.Pointer) error {
	var cErr *C.char
	var ret = C.pluginClose(handler, &cErr)
	if ret != 0 {
		return errors.New(C.GoString(cErr))
	}
	return nil
}
