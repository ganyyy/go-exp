package patch

import (
	"unsafe"
)

func Patch(src, dst unsafe.Pointer) {
	var toBytes = jmpToGoFn(uintptr(dst))
	copyToLocation(uintptr(src), toBytes)
}

func Backup(src unsafe.Pointer) []byte {
	return backFromLocation(uintptr(src))
}

func Restore(src unsafe.Pointer, data []byte) {
	if len(data) != jmpLen {
		return
	}
	copyToLocation(uintptr(src), data)
}
