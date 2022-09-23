package main

import (
	"fmt"
)

type vector[T any] []T

func (v vector[T]) show() {
	for i, i2 := range v {
		fmt.Println(i, i2)
	}
}

type safeMap[K comparable, V any] map[K]V

func (s safeMap[K, V]) show() {
	for k, v := range s {
		fmt.Println(k, v)
	}
}

func main() {
	var intVector vector[int]
	intVector = make([]int, 10)

	intVector.show()

	var mm safeMap[int, string] = make(map[int]string, 1)
	mm[10] = "12345"

	mm.show()

	var a = "sss"
	var ret = GetOrDefault(a, 100)
	var ret2 = GetOrDefault(a, "123131")

	println(ret, ret2)

	println(Add(myInt(10), myInt(20)))
	println(Add(myInt(10), myInt(20)))
	println(Add(myInt(10), myString("20")))
	println(Add(myString("bbb"), myString("bbb")))
	println(Add2(myString("bbb"), myString("bbb")))
	println(Add2(myInt(10), myString("20")))

	var src = []int{1, 2, 3, 4, 5, 6}
	fmt.Println(Map(src, func(t int) int {
		return t + 10
	}))

	var tmp = pack(1, 2, 3, 4, 5, 6, 7)
	var tmp2 = pack("123", "1234", "1235")
	fmt.Println(tmp, tmp2)
	var s set[int] = make(map[int]struct{})
	var s2 set[string] = make(map[string]struct{})
	s.pack(tmp...)
	s2.pack(tmp2...)
	fmt.Println(s, s2)
	fmt.Println(s.unpack(), s2.unpack())
}
