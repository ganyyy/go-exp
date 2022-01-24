//go:build ignore

package main

import (
	"log"
	"reflect"
	"unsafe"
)

func main() {
	var src = []int{1, 2, 3, 4}
	var header = (*reflect.SliceHeader)(unsafe.Pointer(&src))

	var tmp = src
	tmp = append(tmp[:1], tmp[2:]...)
	var tmpHeader = (*reflect.SliceHeader)(unsafe.Pointer(&tmp))

	log.Printf("%+v, %+v", header.Data, tmpHeader.Data)

	log.Printf("%+v, %+v", src, tmp)
}
