package main

import (
	"sync/atomic"
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

func init() {
	go func() {
		var ticker = time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case t := <-ticker.C:
				// println("nowt:", t.String())
				v.Store(t.UnixNano())
			}
		}

	}()

	var nowt = time.Now()
	v.Store(nowt.UnixNano())
}

var v atomic.Int64

func BenchmarkTimeout(b *testing.B) {
	var work = func() {
		var v int
		for i := 0; i < 1e3; i++ {
			v += i
		}
		_ = v
	}

	_ = work

	time.Sleep(time.Second)

	// var next = time.Now().Add(time.Second)
	b.Run("select", func(b *testing.B) {
		var ch = make(chan int)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				work()
				select {
				case <-ch:
				default:
				}
			}
		})
	})
	b.Run("timer", func(b *testing.B) {
		var v int
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				work()
				// v++
				_ = v > 1000
			}
		})
	})
}
