package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
)

const ID = "FuncName"

func doWithGoroutineId(ctx context.Context, id int, f func(context.Context)) {
	go pprof.Do(
		ctx,
		pprof.Labels(ID, strconv.Itoa(id)),
		f,
	)
}

func worker(ctx context.Context) {
	id, ok := pprof.Label(ctx, ID)
	log.Printf("go id:%v, find:%v", id, ok)
	time.Sleep(time.Second * 5)
}

func TestLabel(t *testing.T) {
	var ctx = context.Background()

	for i := 0; i < 5; i++ {
		doWithGoroutineId(ctx, i, worker)
	}

	for _, prof := range pprof.Profiles() {
		var debug int
		if prof.Name() == "goroutine" {
			debug = 1
		} else if prof.Name() == "threadcreate" {
			debug = 2
		}
		prof.WriteTo(os.Stdout, debug)
	}

	time.Sleep(time.Second * 3)
}

func TestGoroutineLabel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	var ctx = context.Background()
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			var labelCtx = pprof.WithLabels(ctx, pprof.Labels(ID, strconv.Itoa(i)))
			pprof.SetGoroutineLabels(labelCtx)
			worker(labelCtx)
		}(i)
	}
	time.Sleep(1 * time.Second)
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	wg.Wait()
	id, ok := pprof.Label(ctx, ID)
	t.Logf("go id:%v, find:%v", id, ok)
}
