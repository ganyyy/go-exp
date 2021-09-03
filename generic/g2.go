package main

import (
	"fmt"
)

func pack[T any](params ...T) []T {
	return params
}

type Student struct {

}

type Addable interface {
	type Student
	ToString() string
}

type set[T comparable] map[T]struct{}

func (s set[T]) pack(params ...T) {
	for _, p := range params {
		s[p] = struct{}{}
	}
}

func (s set[T]) unpack() []T {
	var ret = make([]T, 0, len(s))
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


func main() {
	var tmp = pack[int](1,2,3,4,5,6,7)
	var tmp2 = pack[string]("123","1234","1235")
	fmt.Println(tmp, tmp2)
	var s set[int] = make(map[int]struct{})
	var s2 set[string] = make(map[string]struct{})
	s.pack(tmp...)
	s2.pack(tmp2...)
	fmt.Println(s, s2)
	fmt.Println(s.unpack(), s2.unpack())
}
