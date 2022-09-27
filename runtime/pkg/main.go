package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func Data() {

}

type Stu struct{}

func (s Stu) Name1() {}

func (s *Stu) Name2() {}

var logName = func(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func main() {

	f := runtime.FuncForPC(reflect.ValueOf(Data).Pointer())
	f.Entry()

	fmt.Println(logName(Data))
	fmt.Println(logName(Stu.Name1))
	fmt.Println(logName((*Stu).Name2))
}
