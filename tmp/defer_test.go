package main

import (
	"log"
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

func TestGotoDefer(t *testing.T) {
	var cnt int
next:
	if cnt > 10 {
		return
	}
	defer func(i int) {
		log.Printf("count defer %v", i)
	}(cnt)

	cnt++
	goto next
}
