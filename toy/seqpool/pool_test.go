package seqpool

import (
	"context"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Module struct {
	count     int
	key       int
	isRunning atomic.Bool
}

// IsStop
func (m *Module) IsStop() bool {
	return false
}

func TestPool(t *testing.T) {
	var ctx = context.Background()
	var createCnt atomic.Int64
	var pool = NewPool(ctx, func(i int) *Module {
		createCnt.Add(1)
		return &Module{key: i}
	})

	var now = time.Now()
	var r = rand.New(rand.NewPCG(uint64(now.UnixMilli()), uint64(now.UnixNano())))
	const N = 10
	const C = 5
	const M = 20
	var total atomic.Int64
	var totalSubmit atomic.Int64
	var wg sync.WaitGroup
	var groupCnt [N]atomic.Int64
	wg.Add(C)
	for range C {
		go func() {
			defer wg.Done()
			for range N * M {
				group := r.IntN(N)
				after := groupCnt[group].Add(1)
				wg.Add(1)
				tts := totalSubmit.Add(1)
				_ = tts
				// t.Logf("key: %d, submit: %d, total: %d", group, after, tts)
				pool.Submit(group, func(key int, m *Module) {
					defer wg.Done()
					if !m.isRunning.CompareAndSwap(false, true) {
						t.Logf("key: %d group is multiple running", key)
						return
					}
					defer func() {
						if !m.isRunning.CompareAndSwap(true, false) {
							t.Logf("key: %d group is multiple running", key)
						}
					}()
					tt := total.Add(1)
					_ = tt
					// t.Logf("key: %d, count: %d, total: %d", key, m.count, tt)
					m.count++
					if m.count != int(after) {
						t.Logf("key: %d, count: %d, after: %d not in seq", key, m.count, after)
					}
					time.Sleep(time.Millisecond * 20 * (1 + time.Duration(r.IntN(3))))
				})
			}
		}()
	}
	wg.Wait()

	pool.Stop()

	t.Logf("total: %d, create: %v, empty poll: %v", total.Load(), createCnt.Load(), pool.emptyPoll.Load())
	require.Equal(t, int64(N*M*C), total.Load())

}
