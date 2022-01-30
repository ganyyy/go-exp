package generic2

import (
	"fmt"
	"strings"
)

//Range 任意类型的切片迭代
//优点: 不需要写重复的代码
func Range[T fmt.Stringer](src []T) {
	var tmp = make([]string, 0, len(src))
	for _, v := range src {
		tmp = append(tmp, v.String())
	}
	fmt.Println(strings.Join(tmp, ","))
}

//Range2 使用any作为约束. 除了输出应该没有其他用处
func Range2[T any](src []T) {
	var tmp = make([]string, 0, len(src))
	for _, v := range src {
		tmp = append(tmp, fmt.Sprintf("%+v", v))
	}
	fmt.Println(strings.Join(tmp, ","))
}

//RangeInterface 结构体切片和接口切片无法直接转换
//缺点: 必须要显示进行切片类型的转换
func RangeInterface(src []fmt.Stringer) {
	var tmp = make([]string, 0, len(src))
	for _, v := range src {
		tmp = append(tmp, v.String())
	}
	fmt.Println(strings.Join(tmp, ","))
}

type Stu struct {
	Name string
	Age  int
}

type Stu2 struct {
	Stu
	Other string
}

func (s *Stu) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Name:%v, Age:%v}", s.Name, s.Age)
}

func (s *Stu2) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Stu:%v, Other:%v}", s.Stu.String(), s.Other)
}

func init() {

	/*
		1. Go 的泛型根据实参类型动态的生成相关的函数实现
	*/

	Range(make([]*Stu, 0))
	Range2(make([]*Stu, 0))
	Range(make([]*Stu2, 0))
	Range2(make([]*Stu2, 0))
	Range(make([]fmt.Stringer, 0))
	Range2(make([]fmt.Stringer, 0))
}
