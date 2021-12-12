//go:build ignore
// +build ignore

package main

import (
	"math"
)

func main() {
	var f = 1.0
	println(f + math.NaN())            // NaN
	println(f + math.Inf(0))           // +Inf
	println(f + math.Inf(1))           // +Inf
	println(math.NaN() + math.Inf(1))  // NaN
	println(math.NaN() - math.NaN())   // NaN
	println(math.NaN() + math.Inf(0))  // NaN
	println(math.Inf(1) + math.Inf(0)) // +Inf
	println(math.Inf(1) - math.Inf(1)) // NaN

	var a, b float32 = 1e20, 1e-20
	println(a * a * b)   // +Inf
	println(a * (a * b)) // +1e20
	println(a * (a - a)) // +0
	println(a*a - a*a)   // NaN

	var one = 1.0
	println(one / 0)  // +Inf
	println(one / -0) // +Inf

}
