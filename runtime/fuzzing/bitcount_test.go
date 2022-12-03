package fuzzing

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitCount(t *testing.T) {

	var cases = []int{
		0, 1, 99, 20, 50, math.MaxInt, math.MaxInt64, math.MinInt64, math.MinInt32,
	}

	for _, v := range cases {
		t.Logf("test %v", v)
		assert.Equal(t, loopCount(v), logCount(v))
		assert.Equal(t, loopCount(v), ifCount(v))
	}
}

func BenchmarkBitCount(b *testing.B) {
	n := rand.Int()
	run := func(name string, f func(int) int) {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f(n)
			}
		})
	}

	run("loop", loopCount)
	run("if", ifCount)
	run("log", logCount)
}

func FuzzBitCount(f *testing.F) {
	var cases = []int{
		1, 2, 3, math.MaxInt64, math.MinInt64, math.MinInt32, math.MaxInt32,
	}

	for _, tc := range cases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, a int) {
		loopC := loopCount(a)
		assert.Equal(t, loopC, logCount(a))
		assert.Equal(t, loopC, ifCount(a))
	})
}
