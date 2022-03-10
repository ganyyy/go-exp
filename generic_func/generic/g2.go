//go:build ignore

package main

import (
	"fmt"
	"strconv"
)

func pack[T any](params ...T) []T {
	return params
}

type Student struct {
}

type Addable interface {
	Student
	ToString() string
}

func GetOrDefault[T any](src interface{}, def T) T {
	if t, ok := src.(T); ok {
		return t
	}
	return def
}

type set[T comparable] map[T]struct{}

func (s set[T]) pack(params ...T) {
	for _, p := range params {
		s[p] = struct{}{}
	}
}

func (s set[T]) unpack() []T {
	var ret = make([]T, 0, len(map[T]struct{}(s)))
	for k := range s {
		ret = append(ret, k)
	}
	return ret
}

func packSet[T comparable](params ...T) set[T] {
	var m set[T] = make(map[T]struct{}, len(params))
	for _, p := range params {
		m[p] = struct{}{}
	}
	return m
}

type myInt int

type Addable2 interface {
	~int | ~string
	String() string
}

func (a myInt) String() string {
	return strconv.Itoa(int(a))
}

type myString string

func (s myString) String() string {
	return string(s)
}

func Map[T any](src []T, opt func(T) T) []T {
	var ret = make([]T, len(src))
	for i, v := range src {
		ret[i] = opt(v)
	}
	return ret
}

func Add[T, T2 Addable2](a T2, b T) string {
	return a.String() + b.String()
}

func Add2(a, b fmt.Stringer) string {
	return a.String() + b.String()
}

func main() {

	var a = "sss"
	var ret = GetOrDefault[int](a, 100)
	var ret2 = GetOrDefault[string](a, "123131")

	println(ret, ret2)

	println(Add(myInt(10), myInt(20)))
	println(Add(myInt(10), myInt(20)))
	println(Add(myInt(10), myString("20")))
	println(Add(myString("bbb"), myString("bbb")))
	println(Add2(myString("bbb"), myString("bbb")))
	println(Add2(myInt(10), myString("20")))

	var src = []int{1, 2, 3, 4, 5, 6}
	fmt.Println(Map[int](src, func(t int) int {
		return t + 10
	}))

	var tmp = pack[int](1, 2, 3, 4, 5, 6, 7)
	var tmp2 = pack[string]("123", "1234", "1235")
	fmt.Println(tmp, tmp2)
	var s set[int] = make(map[int]struct{})
	var s2 set[string] = make(map[string]struct{})
	s.pack(tmp...)
	s2.pack(tmp2...)
	fmt.Println(s, s2)
	fmt.Println(s.unpack(), s2.unpack())
}
