package myreflect

import (
	"log"
	"reflect"
)

func CreateStruct(m map[string]int) reflect.Value {

	var rm = reflect.ValueOf(m)
	var iter = rm.MapRange()
	var idx int
	var fields []reflect.StructField
	for iter.Next() {
		var key = iter.Key()

		var structField reflect.StructField
		structField.Anonymous = false
		var val = key.Interface().(string)
		structField.Name = val
		structField.Type = reflect.TypeOf(int(0))
		structField.Index = []int{idx}

		fields = append(fields, structField)
		idx++
	}

	var t = reflect.StructOf(fields)

	var r = reflect.New(t)
	var v = r.Elem()

	for i := 0; i < v.Type().NumField(); i++ {
		var f = v.Field(i)
		var tt = t.Field(i)
		if f.CanSet() {
			log.Printf("field name :%v", tt.Name)
			f.SetInt(int64(m[tt.Name]))
		}
		log.Printf("%+v", f)
	}

	return r
}
