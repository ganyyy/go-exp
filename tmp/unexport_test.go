package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type A struct {
	name string
	Name string
	age  int
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
	fmt.Println(a)
}

func TestUnexport(t *testing.T) {
	unexportField()
}
