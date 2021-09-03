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

func (s safeMap[K,V]) show() {
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
}
