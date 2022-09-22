package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type NonStruct struct {
	Val [100]int
}

func GetVal() *NonStruct {
	return &NonStruct{}
}

const (
	logHeapArenaBytes = 26
	arenaBaseOffset   = 0xFFFF800000000000
	max               = 0x00007FFFFFFFFFFF
	min               = 0x000000C000000000
	All               = max + min

	heapArenaBytes = 1 << logHeapArenaBytes

	arenaL1Bits = 0
	arenaL2Bits = 22

	arenaL1Shift = arenaL2Bits

	arenaBits = arenaL1Bits + arenaL2Bits
)

type arenaIdx uint

func arenaIndex(p uintptr) arenaIdx {
	return arenaIdx((p - arenaBaseOffset) / heapArenaBytes)
}

func arenaBase(i arenaIdx) uintptr {
	return uintptr(i)*heapArenaBytes + arenaBaseOffset
}

func (i arenaIdx) l1() uint {
	return uint(i) >> arenaL1Shift
}

func (i arenaIdx) l2() uint {
	return uint(i) & (1<<arenaL2Bits - 1)

}

func showAddr() {
	var showInfo = func(p uintptr) {
		var idx = arenaIndex(p)
		fmt.Printf("p:%016X, p-Offset: %016X, idx:%016X l1:%016X, l2:%016X\n", p, p-arenaBaseOffset, idx, idx.l1(), idx.l2())
	}

	// var val = GetVal()

	// showInfo(arenaBaseOffset)

	// showInfo(uintptr(unsafe.Pointer(val)))

	// showInfo(0)

	// showInfo(0x7fc000000000)

	// showInfo(min)

	// showInfo(max)

	m := make([]byte, 64<<20)
	showInfo(uintptr((*reflect.SliceHeader)(unsafe.Pointer(&m)).Data))
	showInfo(uintptr((*reflect.SliceHeader)(unsafe.Pointer(&m)).Data + 64<<20))
}
