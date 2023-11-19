package api

import (
	"sync"
	"testing"
)

func TestOnce(t *testing.T) {

	f := sync.OnceFunc(func() {
		t.Log("once")
	})
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()

	v := sync.OnceValue[int](func() int {
		t.Logf("once value")
		return 1
	})

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			t.Log(v())
			wg.Done()
		}()
	}
	wg.Wait()

}
