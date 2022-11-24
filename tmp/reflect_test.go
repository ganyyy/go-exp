package main

import (
	"reflect"
	"testing"
)

type ReflectStruct struct{}

func (r *ReflectStruct) Name() string {
	return reflect.TypeOf(r).Elem().Name()
}

func TestReflect_NilStructElement(t *testing.T) {
	var r ReflectStruct
	t.Logf("struct:%v", r.Name())

	var rp *ReflectStruct
	t.Logf("pointer:%v", rp.Name())
}
