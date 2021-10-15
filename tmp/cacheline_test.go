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
