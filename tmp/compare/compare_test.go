package main_test

import "testing"

func TestCompareInterface(t *testing.T) {
	var a, b any = 1, 20
	var c, d any = []int(nil), []int(nil)

	t.Log(a == 1)
	t.Log(b == a)
	t.Log(c == d)
}
