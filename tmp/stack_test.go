package main

import (
	"testing"
	"unsafe"
)

func TestStackRunning(t *testing.T) {
	var a int
	t.Logf("a:%v, %p", a, (unsafe.Pointer(&a)))

	for i := 0; i < 10000; i++ {
		var t [100]int
		t[0] += 100
	}

	t.Logf("a:%v, %p", a, (unsafe.Pointer(&a)))
}
