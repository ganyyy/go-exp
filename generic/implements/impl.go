package main

import (
	"fmt"
	"reflect"
)

//go:noinline
func AnyPrint[T any](t T) {
	fmt.Println(t)
}

//go:noinline
func AnyPrint2(t any) {
	rt := reflect.TypeOf(t)
	fmt.Println(rt.Kind(), t)
}

//go:noinline
func AnyString[T interface {
	Id() int
	String() string
}](t T) {
	println(t.String(), t.Id())
}

type Struct1 struct {
	Name string
	Age  int
}

func (s *Struct1) String() string { return fmt.Sprintf("Name: %s, Age: %d", s.Name, s.Age) }
func (s *Struct1) Id() int        { return 0 }

type Struct2 struct {
	Age  int
	Name string
}

func (s *Struct2) String() string { return fmt.Sprintf("Name: %s, Age: %d", s.Name, s.Age) }
func (s *Struct2) Id() int        { return 0 }

type Struct3 struct {
	Name2 string
	Age2  int
}

func (s *Struct3) String() string { return fmt.Sprintf("Name: %s, Age: %d", s.Name2, s.Age2) }
func (s *Struct3) Id() int        { return 0 }

type Struct4 struct {
	Name string
	Age  int
}

func (s Struct4) String() string { return fmt.Sprintf("Name: %s, Age: %d", s.Name, s.Age) }
func (s Struct4) Id() int        { return 0 }

type Struct5 struct {
	Name1 string
	Age   int
}

func (s Struct5) String() string { return fmt.Sprintf("Name: %s, Age: %d", s.Name1, s.Age) }
func (s Struct5) Id() int        { return 0 }

func main() {

	AnyPrint(1)
	AnyPrint("hello")
	AnyPrint([]int{1, 2, 3})
	AnyPrint[*int](nil)
	AnyPrint[*bool](nil)
	AnyPrint[*[]int](nil)
	AnyPrint[*struct {
		_ int
		_ string
	}](nil)
	AnyPrint([]int{})
	AnyString(&Struct1{})
	AnyString(&Struct2{})
	AnyString(&Struct3{})
	AnyString(&Struct4{})
	AnyString(Struct4{})
	AnyString(Struct5{})

	var ch = make(chan int)
	AnyPrint2([]int{1, 2, 3})
	AnyPrint2(ch)

}
