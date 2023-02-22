package api

import (
	"sync/atomic"
	"testing"
)

type Data1 struct {
	done uint32
	pad  [40]byte
}

//go:noinline
func (d *Data1) Add() { atomic.AddUint32(&d.done, 1) }

type Data2 struct {
	pad  [40]byte
	done uint32
}

//go:noinline
func (d *Data2) Add() { atomic.AddUint32(&d.done, 1) }

func Benchmark(b *testing.B) {
	b.Run("data1", func(b *testing.B) {
		var data Data1
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				data.Add()
			}
		})
	})
	b.Run("data2", func(b *testing.B) {
		var data Data2
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				data.Add()
			}
		})
	})
}

func TestDeferDefer(t *testing.T) {
	defer func() {
		t.Log("recover = ", recover())
	}()

	defer func() {
		panic("456 ")
	}()

	defer func() {
		panic("in defer")
	}()

	_ = 10
}
