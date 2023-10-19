package helper

import (
	"log/slog"
	"os"
)

func InitSlog() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}),
	)
}
