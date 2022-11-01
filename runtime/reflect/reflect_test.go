package reflect2

import (
	"fmt"
	"reflect"
	"testing"
)

type Student struct {
	Name string
	Age  int
}

func TestStruct(t *testing.T) {
	// rtype是一个基础属性, 作为一个内嵌字段封装到各个类型中
	var stu Student
	var sv = reflect.ValueOf(stu)
	var st = reflect.TypeOf(stu)
	fmt.Println(sv)
	fmt.Println(st.Field(0))
}
