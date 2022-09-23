package main

import "testing"

var v string

func BenchmarkGetID(b *testing.B) {
	b.Run("1", func(b *testing.B) {
		var a AID
		var t string
		for i := 0; i < b.N; i++ {
			t = GetIDDirect(&a)
		}
		v = t
	})

	b.Run("2", func(b *testing.B) {
		var a AID
		var t string
		for i := 0; i < b.N; i++ {
			t = GetID(&a)
		}
		v = t
	})

	b.Run("3", func(b *testing.B) {
		var a AID
		var t string
		for i := 0; i < b.N; i++ {
			t = GetIDInterface(&a)
		}
		v = t
	})
}
