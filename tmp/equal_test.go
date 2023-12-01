//go:build ignore

package main

type A struct {
	X int
	Y string
}

type B struct {
	X int
	Y string
}

type C A

type D = struct {
	X int
	Y string
}

func _() {
	var a A
	var b B
	var c C
	var d D

	_ = a == b // error: mismatched types A and B
	_ = a == c // error: mismatched types A and C
	_ = a == d // ok: struct types are identical
	_ = b == c // error: mismatched types B and C
	_ = b == d // ok: struct types are identical
	_ = c == d // ok: struct types are identical
}
