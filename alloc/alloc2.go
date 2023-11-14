package main

import (
	"fmt"
	"reflect"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

type PtrStruct1 struct {
	_ [131072]byte
	V *byte
	_ [10]byte
}

type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      uint8   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal  func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata *byte // garbage collection data
}

type PtrStruct2 struct {
	_ [1]int
	V *byte
	_ int
}

type PtrStruct3 struct {
	V *byte
	_ [1]int
	_ int
}

type PtrStruct4 struct {
	_ [1]int
	_ int
	V *byte
}

type iface struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}

func showTypeInfo(t interface{}) {
	return
	var typ = reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Array {
		elem := reflect.ValueOf(t).Elem().Index(0).Addr().Pointer()
		elem2 := reflect.ValueOf(t).Elem().Index(1).Addr().Pointer()
		showInfo(elem)
		showInfo(elem2)
		fmt.Printf("element addr: 0x%x 0x%x\n", elem, elem2)
	}
	var dd = (*rtype)((*iface)(unsafe.Pointer(&typ)).data)
	fmt.Printf("%10s:{ptrdata: %2d, gcdata Val:0b%08b}, size:%v\n", typ.String(), dd.ptrdata, *dd.gcdata, typ.Size())
	_ = dd
}

func showGCInfo() {
	var s = make([]struct {
		_ int
		_ *int
		_ int
	}, 1000)

	showTypeInfo(s)
}

var global interface{}
var global2 interface{}

//go:noinline
func makeSomeObject() any {
	return &[50]struct {
		_ int
		_ *int
		_ int
	}{}
}

func showGCBitMap() {

	addr, err := syscall.Mmap(-1, 0, 1<<20, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)

	if err != nil {
		panic(err)
	}

	fmt.Printf("the addr is: %p\n", unsafe.Pointer(unsafe.SliceData(addr)))

	showTypeInfo(new(PtrStruct1))
	//showTypeInfo(new(PtrStruct2))
	//showTypeInfo(new(PtrStruct3))
	//showTypeInfo(new(PtrStruct4))

	// var d = new([(0x28000 / 2)]struct {
	// 	a int
	// 	b *int
	// })

	var obj = makeSomeObject()

	runtime.GC()
	runtime.GC()
	time.Sleep(time.Second * 5)

	// global = d

	// showTypeInfo(d) // {ptrdata： 8, gcdata Val:0b00000010}, size:24

	// d2 := new([1 << 5]struct {
	// 	_ int
	// 	_ *int
	// })

	// showTypeInfo(d2)

	// global2 = d2

	showTypeInfo(new(struct {
		_ *int
		_ *int
		_ int
	})) // {ptrdata：16, gcdata Val:0b00000011}, size:24

	showTypeInfo(new(struct {
		_ *byte
		_ *byte
		_ int
	})) // {ptrdata：16, gcdata Val:0b00000011}, size:24

	showTypeInfo(new(struct {
		_ *byte
		_ int
		_ *byte
	})) // {ptrdata：24, gcdata Val:0b00000101}, size:24

	showTypeInfo(new(struct {
		_ *byte
		_ bool
		_ *byte
	})) // {ptrdata：24, gcdata Val:0b00000101}, size:24

	showTypeInfo(new(struct {
		_ int
		_ *byte
		_ *byte
	})) // {ptrdata：24, gcdata Val:0b00000110}, size:24

	showTypeInfo(new(struct {
		_ int
		_ byte
		_ byte
	})) // {ptrdata： 0, gcdata Val:0b00000000}, size:16

	showTypeInfo(new([3]struct {
		_ *byte
		_ int
	})) // {ptrdata：40, gcdata Val:0b00010101}, size:48

	runtime.GC()
	runtime.GC()
	time.Sleep(time.Second * 5)

	global = obj
	// showTypeInfo(global)
	// showTypeInfo(global2)

	fmt.Println("GC 3")
	runtime.GC()
	runtime.GC()
	time.Sleep(time.Second * 5)

	fmt.Println("GC 4")
	runtime.GC()
	runtime.GC()
	time.Sleep(time.Second * 5)

	fmt.Println("GC 5")
	runtime.GC()
	runtime.GC()
	time.Sleep(time.Second * 5)
}
