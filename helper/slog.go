package helper

import (
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(InitSlog())
}

func InitSlog() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}),
	)
}
