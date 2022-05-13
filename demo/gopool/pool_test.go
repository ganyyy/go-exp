package gopool_test

import (
	"log"
	"strings"
	"testing"
	"time"

	"ganyyy.com/go-exp/gopool"
)

func TestPool(t *testing.T) {
	var p = gopool.NewPool(10)
	p.Start()
	for i := 0; i < 100; i++ {
		var id = i
		p.Run(func() {
			time.Sleep(time.Second)
			log.Printf("my id:%v", id)
		})
	}
	p.Close()
}

func TestReplace(t *testing.T) {
	var replace = func(src, obj, rep string) int {
		var cnt int
		for strings.Contains(src, obj) {
			src = strings.Replace(src, obj, rep, 1)
			cnt++
		}
		return cnt
	}

	t.Logf("cnt:%v", replace("aaaaa", " ", " "))
}
