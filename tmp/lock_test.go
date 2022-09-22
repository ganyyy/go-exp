package main

import (
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	var lock sync.RWMutex

	var start = time.Now()
	var logf = func(info string) {
		t.Logf("nowt %v, %s", time.Since(start), info)
	}

	go func() {
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 300)
		logf("r l 1")
	}()

	go func() {
		time.Sleep(time.Millisecond * 50)
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 400)
		logf("r l 2")
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		lock.Lock()
		defer lock.Unlock()
		time.Sleep(time.Millisecond * 200)
		logf("l 1")
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 100)
		logf("r l 3")
	}()

	time.Sleep(time.Second * 1)
}
