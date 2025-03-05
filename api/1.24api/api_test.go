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
	var m = make(map[int]int)
	_ = m

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
