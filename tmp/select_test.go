package main

import (
	"testing"
	"time"
)

func TestSelectChannel(t *testing.T) {
	var ch = make(chan int, 10)
	var quit = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case v, ok := <-ch:
				t.Logf("v:%v, ok:%v, ln:%v", v, ok, len(ch))
				time.Sleep(time.Millisecond * 500)
			case <-quit:
				t.Logf("close")
				// return
			}
		}
	}()

	for i := 0; i < 10; i++ {
		ch <- i
	}

	time.Sleep(time.Second * 2)
	close(ch)
	time.Sleep(time.Second)
	close(quit)
	time.Sleep(time.Second * 3)

}
