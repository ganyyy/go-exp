package main

import (
	"errors"
	"reflect"
	"syscall"
	"unsafe"
)

type funcVal struct {
	_   uintptr
	ptr unsafe.Pointer
}

func ReplaceFunc(src interface{}, dest interface{}) error {
	var srcType, destType = reflect.TypeOf(src), reflect.TypeOf(dest)
	if srcType.Kind() != destType.Kind() || srcType.Kind() != reflect.Func {
		return errors.New("only replace func")
	}
	// TODO
	return nil
}

func GetFuncAddr(v interface{}) uintptr {
	var value = reflect.ValueOf(v)
	var toAddr = (*funcVal)(unsafe.Pointer(&value)).ptr
	return uintptr(toAddr)
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
	var dest []byte
	var header = (*reflect.SliceHeader)(unsafe.Pointer(&dest))
	header.Data = p
	header.Len = length
	header.Cap = length
	return dest
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
