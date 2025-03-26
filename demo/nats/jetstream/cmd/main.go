package main

import (
	"context"
	"flag"
	"jetstream-demo/common"
	"log/slog"
	"os"
	"sync"
)

func main() {
	flag.Parse()

	common.InitLogger()

	closeCB, err := common.InitNatsConnection()
	if err != nil {
		slog.Default().Error("Failed to connect to NATS", slog.Any("error", err))
		os.Exit(1)
	}
	numProcess := *common.FlagNumProcess
	if numProcess < 1 {
		numProcess = 1
	}
	isProducer := *common.FlagIsProducer
	var runFunc func(context.Context)
	if isProducer {
		runFunc = common.RunProducer
	} else {
		runFunc = common.RunConsumer2
	}

	var wg sync.WaitGroup
	wg.Add(numProcess)

	var ctx, cancel = common.SetupSignalHandler()
	defer cancel()

	defer wg.Wait()
	defer closeCB()

	for idx := range numProcess {
		go func(idx int, ctx context.Context) {
			defer wg.Done()
			runFunc(common.WithRunnerIdx(ctx, idx))
		}(idx, ctx)
	}

	<-ctx.Done()
}
