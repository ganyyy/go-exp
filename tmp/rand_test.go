package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRandNano(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			for j := 0; j < 5; j++ {
				println(time.Now().UnixNano())
			}
		}()
	}

	wait.Wait()

	for i := 0; i < 10; i++ {
		rand.Seed(1)
		t.Logf("%v", rand.Int())
	}

}
