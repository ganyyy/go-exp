package main

import "testing"

func BenchmarkName(b *testing.B) {
	var d = &iData{}
	b.Run("inline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.SetA(i)
		}
	})
	_ = d

	b.Run("no inline", func(b *testing.B) {
		var d2 = &iData{}
		for i := 0; i < b.N; i++ {
			d2.SetB(i)
		}
	})
}
