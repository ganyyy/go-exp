package api

import (
	"fmt"
	"testing"
)

func intRange(n int) int {
	var ret int
	for i := range n {
		ret += i
	}
	return ret
}

func TestRange(t *testing.T) {
	t.Log(intRange(-1))
	t.Log(intRange(10))

	var numbers = make([]int, 10)
	for i := range len(numbers) {
		numbers[i] = i
	}

	for idx, val := range rangeAnyThing(numbers) {
		fmt.Printf("idx: %d, val: %d\n", idx, val)
	}

	for val := range rangeAnyValue(numbers) {
		fmt.Printf("val: %d\n", val)
	}

	for start, end := range groupBatch(100, 7) {
		fmt.Printf("start: %d, end: %d\n", start, end)
	}

	for sub := range groupSlices(numbers, 3) {
		fmt.Printf("sub: %v\n", sub)
	}
}

func groupBatch(total, batch int) func(func(int, int) bool) {
	return func(f func(int, int) bool) {
		for i := 0; i < total; i += batch {
			if !f(i, min(i+batch, total)) {
				break
			}
		}
	}
}

func groupSlices[T any, S ~[]T](s S, batch int) func(func(subSlice S) bool) {
	return func(f func(subSlice S) bool) {
		for i := 0; i < len(s); i += batch {
			if !f(s[i:min(i+batch, len(s))]) {
				break
			}
		}
	}
}

func rangeAnyThing[T any, S ~[]T](s S) func(func(int, T) bool) {
	return func(f func(int, T) bool) {
		for i, v := range s {
			if !f(i, v) {
				break
			}
		}
	}
}

func rangeAnyValue[T any, S ~[]T](s S) func(func(T) bool) {
	return func(f func(T) bool) {
		for _, v := range s {
			if !f(v) {
				break
			}
		}
	}
}
