package seqpool

import (
	"context"
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"time"
)

type Module struct {
	count int
	key   int
}

// IsStop
func (m *Module) IsStop() bool {
	return false
}

func TestPool(t *testing.T) {
	var ctx = context.Background()
	var pool = NewPool[int, *Module](ctx, func(i int) *Module {
		return &Module{key: i}
	})

	var now = time.Now()
	var r = rand.New(rand.NewPCG(uint64(now.UnixMilli()), uint64(now.UnixNano())))
	const N = 5
	var total atomic.Int64
	for i := 0; i < N*10; i++ {
		pool.Submit(r.IntN(N), func(key int, m *Module) {
			t.Logf("key: %d, count: %d", key, m.count)
			m.count++
			time.Sleep(time.Millisecond * 50 * (1 + time.Duration(r.IntN(5))))
			total.Add(1)
		})
	}

	time.Sleep(time.Second * 5)

	pool.Stop()

	t.Logf("total: %d", total.Load())
}
