package common

import (
	"log/slog"
	"os"
)

func InitLogger() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
			}),
	)
	slog.SetDefault(logger)
}
