package timelock

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTimeLock(t *testing.T) {
	var tl = NewTimeLock()

	var wg sync.WaitGroup

	// 协程1申请普通锁, 先执行
	// 协程2申请普通锁, 等待协程1释放普通锁
	// 协程3申请优先级锁, 等待协程1释放普通锁, 并先于协程2执行

	warp := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	warp(func() {
		tl.Lock()
		defer tl.Unlock()
		time.Sleep(3 * time.Second)
		t.Logf("goroutine 1 exit")
	})

	warp(func() {
		time.Sleep(500 * time.Millisecond)
		t.Logf("goroutine 2 try to lock")
		tl.Lock()
		defer tl.Unlock()
		t.Logf("goroutine 2 exit")
	})

	warp(func() {
		time.Sleep(1 * time.Second)
		t.Logf("goroutine 3 try to lock")
		tl.PriorityLock(context.Background())
		defer tl.UnLockPriority()
		t.Logf("goroutine 3 exit")
	})

	for i := 0; i < 10; i++ {
		i := i
		warp(func() {
			time.Sleep(2 * time.Second)
			err := tl.TryLockUntil(time.Second * 2)
			if err != nil {
				t.Logf("goroutine %d lock failed: %v", i+4, err)
				return
			}
			defer tl.Unlock()
			t.Logf("goroutine %d lock success", i+4)
			time.Sleep(time.Second * 2)
		})
	}

	time.Sleep(10 * time.Second)

	for i := 0; i < 10; i++ {
		i := i
		warp(func() {
			time.Sleep(2 * time.Second)
			err := tl.TryLockUntil(time.Second * 2)
			if err != nil {
				t.Logf("goroutine %d lock failed: %v", i+100, err)
				return
			}
			defer func() {
				_ = tl.PriorityLock(context.Background())
				defer tl.UnLockPriority()
				t.Logf("goroutine %d unlock success and lock priority", i+100)
				time.Sleep(time.Second * 1)
			}()
			defer tl.Unlock()
			t.Logf("goroutine %d lock success", i+100)
			time.Sleep(time.Second * 2)
		})
	}

	wg.Wait()
}

func BenchmarkTimeLock(b *testing.B) {
	b.Run("TimeLock", func(b *testing.B) {
		var tl = NewTimeLock()
		var i int
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				tl.Lock()
				i++
				tl.Unlock()
			}
		})
	})

	b.Run("TimeLockPriority", func(b *testing.B) {
		var tl = NewTimeLock()
		var i int
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				tl.PriorityLock(context.Background())
				i++
				tl.UnLockPriority()
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		var mu sync.Mutex
		var i int
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				i++
				mu.Unlock()
			}
		})
	})
}
