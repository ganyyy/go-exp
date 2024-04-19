package common

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalCh
		slog.Default().Info("Received signal", slog.Any("signal", sig))
		cancel()
	}()
	return ctx, cancel
}
