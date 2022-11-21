package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRandNano(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			for j := 0; j < 5; j++ {
				println(time.Now().UnixNano())
			}
		}()
	}

	wait.Wait()

	for i := 0; i < 10; i++ {
		rand.Seed(1)
		t.Logf("%v", rand.Int())
	}

}

func BenchmarkSliceIterator(b *testing.B) {

	var data = make([]int, 100)

	b.ResetTimer()

	var cnt int
	b.Run("Index", func(b *testing.B) {
		data := data
		for i := 0; i < b.N; i++ {
			for i := len(data) - 1; i >= 0; i-- {
				cnt += data[i]
			}
		}
	})
	cnt = 0

	b.Run("Range", func(b *testing.B) {
		data := data
		for i := 0; i < b.N; i++ {
			for _, v := range data {
				cnt += v
			}
		}
	})

	_ = cnt
}
