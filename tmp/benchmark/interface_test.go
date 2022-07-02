package main

import (
	"testing"
)

func BenchmarkDoTask(b *testing.B) {
	var total = len(arr)
	var idx uint64
	var randTask = func() Task {
		idx++
		return arr[idx%uint64(total)]
	}
	b.Run("Interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DoTask(randTask())
		}
	})

	b.Run("Struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DoTask2(randTask())
		}
	})
}
