package reflect2

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type Student struct {
	Name string
	Age  int
}

type student struct {
	Name string
	Age  int
}

type privateData struct {
	name string
	age  int
}

var ss = struct {
	Name string
	Age  int
}{}

func TestStruct(t *testing.T) {
	// rtype是一个基础属性, 作为一个内嵌字段封装到各个类型中
	var stu Student
	var sv = reflect.ValueOf(stu)
	var st = reflect.TypeOf(stu)
	fmt.Println(sv)
	fmt.Println(st.Field(0))
}

func TestReflectString(t *testing.T) {
	printStruct := func(v interface{}) {
		vt := reflect.TypeOf(v)
		t.Logf("%s.%s, %s", vt.PkgPath(), vt.Name(), vt.String())
	}

	printStruct(Student{})
	printStruct(student{})
	printStruct(ss)
}

func TestModifyStruct(t *testing.T) {
	var s = &privateData{}

	rt := reflect.TypeOf(s).Elem()

	nameType, ok := rt.FieldByName("name")
	require.True(t, ok)
	*(*string)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + nameType.Offset)) = "hello"

	t.Logf("%+v", *s)
}

func _() {
	var data struct {
		A1 int
		B1 *int
		C1 *struct {
			A2 int
			B2 *int
		}
		D1 *struct {
			A3 int
			B3 int
		}
	}

	_ = data
}

/*
	当data被标记为存活时, 意味着data对应的这块内存整体都是存活的.
	A1本身就是data的一部分, 不需要额外标记.
	但是, data中存在一些指针类型字段, 这些指针类型指向的内存区域, 不一定是data的一部分,
	B1, C1, D1都是指针类型, 如果指向的是合法的内存区域, 这些内存会被标记为存活.

	data的内存布局如下: 包括其引用的对象的内存布局.
	+----+
	| A1 |
	+----+
	|*B1 | --------------------------------> +----+
	+----+									 |	  |
	|*C1 | ------------------->	+----+		 +----+
	+----+				   		| A2 |
	|*D1 | ---> +----+			+----+
	+----+      | A3 |	   		|*B2 | ----> +----+
				+----+			+----+	     |	  |
				| B3 |						 +----+
				+----+

	C1和B1, D1的不同之处在于: C1指向的内存区域是存在指针类型的, B1和D1指向的内存区域是不存在指针的.
	所以, C1指向的内存区域需要进一步的扫描(也就是*B2), 而B1和D1指向的内存区域不需要进一步扫描.
*/
