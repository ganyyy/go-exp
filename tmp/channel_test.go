package main

import "testing"

func TestSlice(t *testing.T) {
	a := make([]int, 0, 2)
	t.Log(a)
	b := append(a, 1)
	t.Log(b)
	c := append(a, 2)
	t.Log(c)

	t.Log(a, b, c)
}
