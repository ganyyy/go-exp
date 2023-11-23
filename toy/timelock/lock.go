package timelock

import (
	"context"
	"time"
)

/*
TimeLock: 提供普通锁和优先级锁
获取锁时, 优先级锁优先级高于普通锁
实际控制数据的锁是优先级锁, 普通锁只是用来限制普通申请锁的数量, 并且优先级锁的申请不受普通锁的限制
这使得优先级锁可以在普通锁被占用时, 优先获取锁
*/
type TimeLock struct {
	sema     chan struct{} // 普通chan,
	priority chan struct{} // 优先级chan,
}

// NewTimeLock 创建一个TimeLock
func NewTimeLock() *TimeLock {
	return &TimeLock{
		sema:     make(chan struct{}, 1),
		priority: make(chan struct{}, 1),
	}
}

// PriorityLock: 优先级锁
func (tl *TimeLock) PriorityLock(ctx context.Context) error {
	return tl.lockWithPriority(true, ctx)
}

// lock 申请普通锁
func (tl *TimeLock) lock(ctx context.Context) error {
	return tl.lockWithPriority(false, ctx)
}

// lockPriority 申请优先级锁
func (tl *TimeLock) lockPriority(ctx context.Context) error {
	select {
	case tl.priority <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// lockWithPriority 申请锁
func (tl *TimeLock) lockWithPriority(priority bool, ctx context.Context) error {
	if priority {
		return tl.lockPriority(ctx)
	} else {
		select {
		case tl.sema <- struct{}{}:
			err := tl.lockPriority(ctx)
			if err != nil {
				// 如果获取优先级锁失败，需要释放普通锁
				<-tl.sema
			}
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// unlockWithPriority 释放锁
func (tl *TimeLock) unlockWithPriority(priority bool) {
	select {
	case <-tl.priority:
	default:
		panic("unlock of unlocked mutex priority")
	}
	if priority {
		return
	}
	select {
	case <-tl.sema:
	default:
		panic("unlock of unlocked mutex sema")
	}
}

// Unlock 普通锁
func (tl *TimeLock) Unlock() {
	tl.unlockWithPriority(false)
}

// UnLockPriority 优先级锁
func (tl *TimeLock) UnLockPriority() {
	tl.unlockWithPriority(true)
}

// Lock: 一直阻塞，直到获取到锁
func (tl *TimeLock) Lock() {
	_ = tl.lock(context.Background())
}

// TryLockUntil: 尝试获取锁，直到超时
func (tl *TimeLock) TryLockUntil(duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return tl.lock(ctx)
}
