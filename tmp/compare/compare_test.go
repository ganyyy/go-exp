package main_test

import "testing"

func TestCompareInterface(t *testing.T) {
	var a, b any = 1, 20
	var _, _ any = []int(nil), []int(nil)

	t.Log(a == b)
	// t.Log(c == d)
}
