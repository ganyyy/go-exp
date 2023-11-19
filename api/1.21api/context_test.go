package api

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	var wg sync.WaitGroup
	logCtx := func(ctx context.Context) {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			<-ctx.Done()
			err := ctx.Err()
			t.Log(err, "err:", context.Cause(ctx))
		}(ctx)
	}

	{
		// 普通的取消
		ctx, cancel :=
			context.WithCancelCause(context.Background())
		logCtx(ctx)
		cancel(errors.New("cancel"))
	}

	{
		// 超时但是手动取消
		ctx, cancel :=
			context.WithTimeoutCause(context.Background(), time.Second, errors.New("timeout"))
		logCtx(ctx)
		cancel()
	}

	{
		// 超时
		ctx, _ :=
			context.WithTimeoutCause(context.Background(), time.Second, errors.New("timeout"))
		logCtx(ctx)
	}

	wg.Wait()
}

func TestContext2(t *testing.T) {
	// 这个顺序很有意思
	// 如果是父节点的cancel在前, 那么注册的after func会被执行
	// 如果是子节点自己先cancel, 那么注册的after func不会被执行
	{
		ctx, cancel := context.WithCancel(context.Background())
		fn := context.AfterFunc(ctx, func() {
			t.Log("after func")
		})
		// 这个顺序会执行after func
		cancel()
		assert.False(t, fn())
	}

	// 这个顺序不会执行after func
	{
		ctx, cancel := context.WithCancel(context.Background())
		fn := context.AfterFunc(ctx, func() {
			t.Log("after func")
		})
		assert.True(t, fn())
		cancel()
	}
	time.Sleep(time.Second)
}
