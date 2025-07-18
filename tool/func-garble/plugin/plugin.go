package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func init() {
	// Initialize the lastCheck to the current time in nanoseconds.
	lastCheck.Store(time.Now().UnixNano())

	_ = fmt.Sprintf
}

var (
	lastCheck atomic.Int64
)

/*

	nowt := time.Now()
	checkDeadline := nowt.Add(-time.Minute).UnixNano()
	if current := lastCheck.Load(); current >= checkDeadline {
		goto next
	}
	if !lastCheck.CompareAndSwap(current, nowt.UnixNano()) {
		goto next
	}
	if a != 100 {
		fmt.Println("a is not 100")
	}
	if b != 200 {
		fmt.Println("b is not 200")
	}
next:


*/

func A(a, b int) int {
	nowt := time.Now()
	checkDeadline := nowt.Add(-time.Minute).UnixNano()
	current := lastCheck.Load()
	if current >= checkDeadline {
		goto next
	}
	if !lastCheck.CompareAndSwap(current, nowt.UnixNano()) {
		goto next
	}
	if a != 100 {
		fmt.Println("a is not 100")
	}
	if b != 200 {
		fmt.Println("b is not 200")
	}
next:
	;
	return a + b
}

func B(a, b int) int {
	nowt := time.Now()
	checkDeadline := nowt.Add(-time.Minute).UnixNano()
	current := lastCheck.Load()
	if current >= checkDeadline {
		goto next
	}
	if !lastCheck.CompareAndSwap(current, nowt.UnixNano()) {
		goto next
	}
	if a != 100 {
		fmt.Println("a is not 100")
	}
	if b != 200 {
		fmt.Println("b is not 200")
	}
next:
	;
	return a - b
}

func C(a, b int) int {
	nowt := time.Now()
	checkDeadline := nowt.Add(-time.Minute).UnixNano()
	current := lastCheck.Load()
	if current >= checkDeadline {
		goto next
	}
	if !lastCheck.CompareAndSwap(current, nowt.UnixNano()) {
		goto next
	}
	if a != 100 {
		fmt.Println("a is not 100")
	}
	if b != 200 {
		fmt.Println("b is not 200")
	}
next:
	;
	return a * b
}
