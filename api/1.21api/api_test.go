package api

import (
	"context"
	"math"
	"os"
	"runtime"
	"strings"
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

type MyHandler struct {
	log slog.Handler
	t   *testing.T
}

// Enabled implements slog.Handler.
func (m *MyHandler) Enabled(c context.Context, l slog.Level) bool {
	return m.log.Enabled(c, l)
}

// WithAttrs implements slog.Handler.
func (m *MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MyHandler{
		log: m.log.WithAttrs(attrs),
		t:   m.t,
	}
}

// WithGroup implements slog.Handler.
func (m *MyHandler) WithGroup(name string) slog.Handler {
	return &MyHandler{
		log: m.log.WithGroup(name),
		t:   m.t,
	}
}

func (m *MyHandler) Handle(ctx context.Context, r slog.Record) error {
	// TODO 如果要实现边长的stack?
	if r.Level < slog.LevelError {
		return m.log.Handle(ctx, r)
	}
	// 更好的实现, 可以参考zap的实现, 这里只是简单的展示一个例子
	var pcs [16]uintptr
	n := runtime.Callers(4, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	var stack []string
	for {
		frame, more := frames.Next()
		stack = append(stack, frame.Function)
		if !more {
			break
		}
	}
	r.AddAttrs(slog.String("stack", strings.Join(stack, "\n")))
	return m.log.Handle(ctx, r)
}

func dfsCall(i int) int {
	if i < 10 {
		slog.Debug("dfs call", slog.Int("i", i))
		return dfsCall(i + 1)
	}
	slog.Error("dfs call end!", slog.Int("i", i))
	return i
}

func TestSLog(t *testing.T) {
	defer func() { _ = os.Stdout.Sync() }()
	var l = slog.New(&MyHandler{
		t: t,
		log: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	})

	l = l.WithGroup("test").With(slog.String("name", "123"), slog.Int("age", 10))
	l.Debug("group test", slog.String("name", "456"))
	slog.SetDefault(l)
	dfsCall(0)
}
