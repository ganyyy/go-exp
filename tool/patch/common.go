package patch

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
	if len(str) == 0 {
		return (*C.char)(NULL)
	}
	var name = append([]byte(str), 0) // 增加一个\0终止符
	return (*C.char)(unsafe.Pointer(&name[0]))
}

var NULL = unsafe.Pointer((*int)(nil))

func FuncAddr(fun interface{}) unsafe.Pointer {
	var value = reflect.ValueOf(fun)
	return (*funcVal)(unsafe.Pointer(&value)).ptr
}
