package api

import (
	"context"
	"math"
	"os"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	require.Equal(t, 5, max(1, 2, 3, 4, 5))
	require.Equal(t, 1, min(1, 2, 3, 4, 5))
	require.Equal(t, "1", min("1", "2", "3", "4", "5"))

	var _ = math.MaxInt

	var s = []int{1, 2, 3, 4, 5}
	clear(s)
	require.Equal(t, make([]int, 5), s)
}

func TestSLog(t *testing.T) {
	defer func() {
		_ = os.Stdout.Sync()
	}()
	var l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			_ = groups
			return a
		},
	}))

	l.Error("hello %s, %d, %s", "world", 123, "abc")

	l.LogAttrs(context.Background(), slog.LevelDebug, "hello", slog.String("name", "123"), slog.Int("age", 10))
	slog.SetDefault(l)
	slog.Debug("hello %s, %d, %s", "world", 123, "abc")
}
