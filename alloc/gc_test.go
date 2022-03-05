package main

import (
	"runtime"
	"testing"
	"time"
)

func TestGC(t *testing.T) {
	t.Run("before", func(t *testing.T) {
		type Before struct {
			a [10]*byte
			_ [1 << 15]int
		}

		var arr = make([]Before, 1<<10)

		for i := 0; i < 10; i++ {
			var now = time.Now()
			runtime.GC()
			t.Logf("cost %v", time.Now().Sub(now))
		}

		runtime.KeepAlive(arr)
	})

	t.Run("after", func(t *testing.T) {
		type After struct {
			_ [1 << 15]int
			a [10]*byte
		}

		var arr = make([]After, 1<<10)

		for i := 0; i < 10; i++ {
			var now = time.Now()
			runtime.GC()
			t.Logf("cost %v", time.Now().Sub(now))
		}

		runtime.KeepAlive(arr)
	})
}
