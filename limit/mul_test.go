package main

import "testing"

type T = int32

func op(a, b T) T {
	return a + b
}

func BenchmarkMulti(b *testing.B) {

	var src = make([]T, 1e4)
	for i := range src {
		src[i] = T(i + 1)
	}
	var ln = len(src)

	b.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ret T = 1
			for j := 0; j < ln; j++ {
				ret = op(ret, src[j])
			}
			_ = ret
		}
	})

	b.Run("Op", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ret T = 1
			var t1, t2 T = 1, 1
			for j := 0; j < ln; j += 2 {
				t1 = op(t1, src[j])
				t2 = op(t2, src[j+1])
			}
			ret = t1 * t2
			_ = ret
		}
	})

	b.Run("loop1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ret T = 1
			for j := 0; j < ln; j += 2 {
				ret = (ret * src[j]) * src[j+1]
			}
			_ = ret
		}
	})

	b.Run("Op2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ret T = 1
			for j := 0; j < ln; j += 2 {
				ret = ret * (src[j] * src[j+1])
			}
			_ = ret
		}
	})
}
