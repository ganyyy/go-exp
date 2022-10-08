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

func TestMethMod(t *testing.T) {
	a, b := 5, -5

	logMod := func(a, b int) {
		t.Logf("%d %% %d = %d", a, b, a%b)
	}

	logMod(a, -3)
	logMod(a, 3)
	logMod(b, -3)
	logMod(b, 3)

	const k = 5

	logMod2 := func(i, id int) {
		v := i + ((id-i)%k+k)%k
		t.Logf("i:%v, id:%v, nextId:%v", i, id, v)
	}

	for i := 0; i < 10; i++ {
		logMod2(i, 3)
	}
}
