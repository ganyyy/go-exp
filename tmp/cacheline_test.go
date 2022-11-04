package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkCacheLineTest(b *testing.B) {
	b.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var an, bn int32
			var wait sync.WaitGroup
			wait.Add(2)
			go func() {
				defer wait.Done()
				for i := 0; i < 10000; i++ {
					atomic.AddInt32(&an, 1)
				}
			}()

			go func() {
				defer wait.Done()
				for i := 0; i < 10000; i++ {
					atomic.AddInt32(&bn, 1)
				}
			}()

			wait.Wait()
		}
	})

	b.Run("Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var an, bn [16]int32
			var wait sync.WaitGroup
			wait.Add(2)
			go func() {
				defer wait.Done()
				for i := 0; i < 10000; i++ {
					atomic.AddInt32(&an[0], 1)
				}
			}()

			go func() {
				defer wait.Done()
				for i := 0; i < 10000; i++ {
					atomic.AddInt32(&bn[0], 1)
				}
			}()

			wait.Wait()
		}
	})
}

const N = 1000

var a [N]int

//go:noinline
func g0(a *[N]int) {
	for i := range a {
		a[i] = i // line 12
	}
}

//go:noinline
func g1(a *[N]int) {
	_ = *a // line 18
	for i := range a {
		a[i] = i // line 20
	}
}

func BenchmarkGGG(b *testing.B) {
	b.Run("g0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g0(&a)
		}
	})
	b.Run("g1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g1(&a)
		}
	})
}
