package main

import (
	"testing"

	"ganyyy.com/go-exp/demo/hotfix/common"
)

func BenchmarkName(b *testing.B) {
	var d = &common.Data{}
	b.Run("inline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.SetA(i)
		}
	})
	_ = d

	b.Run("no inline", func(b *testing.B) {
		var d2 = &common.Data{}
		for i := 0; i < b.N; i++ {
			d2.SetB(i)
		}
	})
}
