package api

import (
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"
	"weak"
)

func TestMap(t *testing.T) {
	var m = make(map[int]int, 8)
	_ = m

	for i := 1; i <= 8; i++ {
		m[i] = i
	}

	for k, v := range m {
		for i := 1; i <= 10; i++ {
			m[k*100+i] = i
		}
		fmt.Println(k, v)
	}
	fmt.Println(len(m))
}

func TestMap2(t *testing.T) {
	t.Run("NaNKey", func(t *testing.T) {
		var m = make(map[float64]int, 8)
		m[math.NaN()] = 1
		m[math.NaN()] = 2
		m[math.NaN()] = 3
		fmt.Println(m)
		for k, v := range m {
			fmt.Println(k, v)
		}
	})
}

func TestWeak(t *testing.T) {
	var val = new(int)
	var v int
	runtime.AddCleanup(val, func(v int) {
		fmt.Println("cleanup", v)
	}, v)
	ptr := weak.Make(val)
	runtime.GC()
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)

	pv := ptr.Value()
	if pv != nil {
		panic("value should be nil")
	}

}

func TestMap2(t *testing.T) {
	var m = make(map[float64]int)
	m[math.NaN()] = 1
	m[math.NaN()] = 2
	m[math.NaN()] = 3
	m[math.Inf(1)] = 2
	m[math.Inf(-1)] = 3

	delete(m, math.NaN())
	delete(m, math.Inf(1))
	delete(m, math.Inf(-1))

	t.Logf("map: %v", m)

	for k, v := range m {
		t.Logf("k: %v, v: %v", k, v)
		delete(m, k)
	}

	t.Logf("map: %v", m)
}

func BenchmarkName(b *testing.B) {

	for b.Loop() {

	}
}

// func TestSyncTest(t *testing.T) {
// 	synctest.Run(func() {
// 		fmt.Println("hello")
// 		synctest.Wait()
// 		fmt.Println("world")
// 	})
// }
