package patch

/*
#include <stdlib.h>
*/
import "C"

import (
	"reflect"
	"unsafe"
)

type funcVal struct {
	_   uintptr
	ptr unsafe.Pointer
}

func toCString(str string) *C.char {
	return C.CString(str)
}

func freeCString(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

var NULL = unsafe.Pointer((*int)(nil))

func FuncAddr(fun interface{}) unsafe.Pointer {
	var value = reflect.ValueOf(fun)
	// unsafe.Pointer(runtime.FuncForPC(value.Pointer()).Entry())
	return (*funcVal)(unsafe.Pointer(&value)).ptr
}
