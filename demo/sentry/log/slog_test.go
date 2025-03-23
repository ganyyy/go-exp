package log

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

var panicExtraInfo struct{}

func init() {
	sentry.AddGlobalEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {

		if event == nil || event.Level != sentry.LevelFatal || hint == nil || hint.Context == nil {
			return event
		}

		val := hint.Context.Value(&panicExtraInfo)
		if val == nil {
			return event
		}
		// 这是一个通过 sentry.RecoverWithContext 捕获的 panic
		msgGen, ok := val.(func() string)
		if ok && msgGen != nil {
			event.Message = msgGen() + "\n" + event.Message
		}
		return event
	})
}

func RuntimeStack(skip int) *sentry.Stacktrace {
	stack := sentry.NewStacktrace()
	if skip <= len(stack.Frames) {
		stack.Frames = stack.Frames[skip:]
	}
	return stack
}

func R(format string, args ...any) context.Context {
	return context.WithValue(context.Background(), &panicExtraInfo, func() string {
		return fmt.Sprintf(format, args...)
	})
}

type PipeLogHandle []slog.Handler

// Enabled implements slog.Handler
func (h PipeLogHandle) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle implements slog.Handler
func (h PipeLogHandle) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range h {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}
		if err := handler.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WithAttrs implements slog.Handler
func (h PipeLogHandle) WithAttrs(attrs []slog.Attr) slog.Handler {
	var handlers = make(PipeLogHandle, 0, len(h))
	for _, handler := range h {
		handlers = append(handlers, handler.WithAttrs(attrs))
	}
	return handlers
}

// WithGroup implements slog.Handler
func (h PipeLogHandle) WithGroup(name string) slog.Handler {
	var handlers = make(PipeLogHandle, 0, len(h))
	for _, handler := range h {
		handlers = append(handlers, handler.WithGroup(name))
	}
	return handlers
}

var LogLevel = &slog.LevelVar{}

func TestLogger(t *testing.T) {

	LogLevel.Set(slog.LevelDebug)

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://6b9a1597deb83afb188c24dcaaa5bd1e@o4509022030856192.ingest.us.sentry.io/4509022051172352",
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(10 * time.Second)

	logger := slog.New(PipeLogHandle{
		sentryslog.Option{
			Level:     slog.LevelError,
			AddSource: true,
			Converter: func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record, hub *sentry.Hub) *sentry.Event {
				event := sentryslog.DefaultConverter(addSource, replaceAttr, loggerAttr, groups, record, hub)
				event.Threads = append(event.Threads, sentry.Thread{
					Current:    true,
					Stacktrace: RuntimeStack(2),
				})
				return event
			},
		}.NewSentryHandler(),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     LogLevel,
			AddSource: true,
		}),
	})

	logger = logger.With(
		slog.String("release", "1.0.0"),
	).With(
		slog.Group("test", slog.String("name", "test"), slog.Int("number", 123)),
	)

	slog.SetDefault(logger)

	manyCall(10)

	panicFn()
}

func manyCall(depth int) {
	if depth == 0 {
		slog.Default().Info("This is a test message", slog.Int("info", 1))
		slog.Default().Error("This is an error message", slog.Int("error", 2))
		slog.Default().Warn("This is a warning message", slog.Int("warn", 3))
		slog.Default().Debug("This is a debug message", slog.Int("debug", 4))
		return
	}
	manyCall(depth - 1)
}

func panicFn() {
	defer sentry.RecoverWithContext(R("test some panic %v", 82))
	panic("test panic")
}

func TestLogger2(t *testing.T) {

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://6b9a1597deb83afb188c24dcaaa5bd1e@o4509022030856192.ingest.us.sentry.io/4509022051172352",
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(10 * time.Second)

	logger := NewLogger(Option{
		MinLogLevel:     slog.LevelDebug,
		StacktraceLevel: slog.LevelWarn,
		LogOption: slog.HandlerOptions{
			AddSource: true,
		},
		SentryLogOption: sentryslog.Option{
			AddSource: true,
		},
	})

	var manyCall func(depth int)

	manyCall = func(depth int) {
		if depth == 0 {
			logger.Info("This is a test message", slog.Int("info", 1))
			logger.Error("This is an error message", slog.Int("error", 2))
			logger.Warn("This is a warning message", slog.Int("warn", 3))
			logger.Debug("This is a debug message", slog.Int("debug", 4))
			return
		}
		manyCall(depth - 1)
	}
	manyCall(5)
}
