package api

import (
	"sync"
	"testing"
)

func TestSync(t *testing.T) {
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

}
