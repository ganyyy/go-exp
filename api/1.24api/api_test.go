package api

import (
	"fmt"
	"runtime"
	"testing"
	"testing/synctest"
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

func BenchmarkName(b *testing.B) {

	for b.Loop() {

	}
}

func TestSyncTest(t *testing.T) {
	synctest.Run(func() {
		fmt.Println("hello")
		synctest.Wait()
		fmt.Println("world")
	})
}
