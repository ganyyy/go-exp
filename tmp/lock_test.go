package main

import (
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var lock sync.RWMutex

	go func() {
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 300)
		t.Logf("r l 1")
	}()

	go func() {
		time.Sleep(time.Millisecond * 50)
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 400)
		t.Logf("r l 2")
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		lock.Lock()
		defer lock.Unlock()
		time.Sleep(time.Millisecond * 200)
		t.Logf("l 1")
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		lock.RLock()
		defer lock.RUnlock()
		time.Sleep(time.Millisecond * 100)
		t.Logf("r l 3")
	}()

	time.Sleep(time.Second * 1)
}
