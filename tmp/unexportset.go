package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type A struct {
	name string
	age  int
	Name string
}

func unexportField() {

	var a A

	var rv = reflect.ValueOf(&a).Elem()
	var rt = reflect.TypeOf(&a).Elem()

	for i := 0; i < rv.NumField(); i++ {
		var field = rv.Field(i)
		var ft = rt.Field(i)
		if !ft.IsExported() {
			if ft.Name == "name" {
				// 不推荐这么搞, 中间可能会被GC
				*(*string)(unsafe.Pointer(field.UnsafeAddr())) = "456789"
			}
			if ft.Name == "age" {
				*(*int)(unsafe.Pointer(field.UnsafeAddr())) = 100
			}
		} else {
			field.SetString("99999")
		}
		fmt.Println(field)
	}
}
