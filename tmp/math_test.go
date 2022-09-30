package main

import (
	"math"
	"testing"
)

func TestMathInf(t *testing.T) {
	var f = 1.0
	t.Log(f + math.NaN())            // NaN
	t.Log(f + math.Inf(0))           // +Inf
	t.Log(f + math.Inf(1))           // +Inf
	t.Log(math.NaN() + math.Inf(1))  // NaN
	t.Log(math.NaN() - math.NaN())   // NaN
	t.Log(math.NaN() + math.Inf(0))  // NaN
	t.Log(math.Inf(1) + math.Inf(0)) // +Inf
	t.Log(math.Inf(1) - math.Inf(1)) // NaN

	var a, b float32 = 1e20, 1e-20
	t.Log(a * a * b)   // +Inf
	t.Log(a * (a * b)) // +1e20
	t.Log(a * (a - a)) // +0
	t.Log(a*a - a*a)   // NaN

	var one = 1.0
	t.Log(one / 0)  // +Inf
	t.Log(one / -0) // +Inf

	t.Log(math.Inf(0) + math.NaN())
	t.Log(math.Inf(-1) + 1)
}
