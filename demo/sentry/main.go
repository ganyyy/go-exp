package main

import (
	"time"

	"ganyyy.com/go-exp/demo/sentry/log"
)

func main() {
	// log.Errorf("Hello World!")

	func() {
		defer log.Recover("panic in %v", "func")
		panic("error")
	}()

	time.Sleep(time.Second * 10)
}
