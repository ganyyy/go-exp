package patch

import (
	"syscall"
	"unsafe"
)

func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
}

func backFromLocation(location uintptr) []byte {
	f := rawMemoryAccess(location, jmpLen)
	var ret = make([]byte, len(f))
	copy(ret, f)
	return ret
}

func pageStart(ptr uintptr) uintptr {
	return ptr &^ (uintptr(syscall.Getpagesize() - 1))
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
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), length)
}

var (
	jmpLen = len(jmpGo)
)

func littleEndian(to uintptr) []byte {
	return []byte{
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56),
	}
}

func jmpToGoFn(to uintptr) []byte {
	var tmp = make([]byte, len(jmpGo))
	copy(tmp, jmpGo[:])
	copy(tmp[2:], littleEndian(to))
	return tmp
}
