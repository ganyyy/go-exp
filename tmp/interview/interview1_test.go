package main

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
)

func TestGetStructTag(t *testing.T) {
	type MyStruct struct {
		A int `tag:"a"`
		B int `tag:"b"`
		c int `tag:"c"`
	}

	var a MyStruct

	t.Logf("%+v", a)

	var rv = reflect.ValueOf(&a)
	var rt = rv.Type()
	for i := 0; i < rv.Elem().NumField(); i++ {
		var field = rv.Elem().Field(i)
		var fieldType = rt.Elem().Field(i)
		t.Logf("field %v tag %v, can set:%v", fieldType.Name, fieldType.Tag.Get("tag"), field.CanSet())
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Int:
				field.SetInt(100 + int64(i))
			}
		}
	}
	t.Logf("%+v", a)

}

func TestZeroSlice(t *testing.T) {
	var a []int
	var b = make([]int, 0)
	var c = make([]int, 0, 0)

	var bs []byte

	bs, _ = json.Marshal(a)
	t.Logf("%+v", string(bs))

	bs, _ = json.Marshal(b)
	t.Logf("%+v", string(bs))

	bs, _ = json.Marshal(c)
	t.Logf("%+v", string(bs))
}

func TestNilSelect(t *testing.T) {
	var c chan int

	var a, b interface{} = 10, []int{10}

	t.Log(a == b)

	select {
	case <-c:
	default:
		t.Logf("default case")
	}

	t.Logf("done")
}

func TestSyncMap(t *testing.T) {
	var m sync.Map
	m.Store("123", "1323")
}
