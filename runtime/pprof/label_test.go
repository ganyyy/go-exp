package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"
)

const ID = "GoId"

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
	time.Sleep(time.Minute * 5)
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
