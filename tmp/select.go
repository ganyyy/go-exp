//go:build ignore

package main

func SelectCase() {
	var ch1 = make(chan int)
	var ch2 = make(chan int)
	var ch3 = make(chan int)
	var ch4 = make(chan int)
	var ch5 = make(chan int)

	select {
	case <-ch1:
	case ch2 <- 2:
	case <-ch3:
	case ch4 <- 1:
	case <-ch5:
	default:
	}
}
