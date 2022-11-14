package main

import (
	"testing"
	"time"
)

func TestTimeAfterPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("recover error:%v", err)
		}
	}()

	var closure = func() {
		panic("panic in time after")
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Logf("recover error:%v", err)
			}
		}()
		closure()
	}()

	time.AfterFunc(time.Second, closure)

	time.Sleep(time.Second * 2)
}
