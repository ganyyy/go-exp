package main

import (
	"fmt"
)

func main() {

	var a int

	var updateA = func() int {
		a++
		return a
	}

	var t = struct {
		A, B, C, D int
	}{
		B: updateA(),
		C: updateA(),
		A: updateA(),
		D: a,
	}

	fmt.Printf("%+v", t) // {A:3 B:1 C:2 D:3}

}
