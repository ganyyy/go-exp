package package_b

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

type inner struct {
	name string
	age  int
	tmp  [1 << 20]int
}

func (i *inner) show() {
	log.Printf("%+v", *i)
}

func GetInner(name string, age int) *inner {
	var in = &inner{
		name: name,
		age:  age,
	}
	in.tmp[0] = 100
	return in
}

type Show interface {
	show()
}

type Outer struct {
	Name string
	Age  int
}

type FuncVal struct {
	_   uintptr
	Ptr unsafe.Pointer
}

func Add(a, b int) int {
	fmt.Printf("Call Add:")
	return a + b
}

func ShowFuncAddr(val interface{}) {
	var value = reflect.ValueOf(val)
	if value.Kind() != reflect.Func {
		fmt.Println("Must input func!")
		return
	}

	fmt.Printf("the func addr is:%v\n", (*FuncVal)(unsafe.Pointer(&val)).Ptr)
}

func ShowAddAddr() {
	ShowFuncAddr(Add)
}

func NewStu() *stu {
	return &stu{}
}

type stu struct {
	name string
}

func (s *stu) setName(newName string) {
	s.name = newName
}
