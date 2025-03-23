package log

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

type Logger struct {
	stacktraceLevel slog.Leveler
	consoleLog      slog.Handler
	sentryLog       slog.Handler
}

// Enabled implements slog.Handler
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.consoleLog.Enabled(ctx, level) || l.sentryLog.Enabled(ctx, level)
}

var StacktraceKey = "stacktrace"

// Handle implements slog.Handler
func (l *Logger) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	if l.stacktraceLevel.Level() <= record.Level {
		stacktrace := sentry.NewStacktrace()
		ln := len(stacktrace.Frames)
		if ln >= 3 {
			// 去掉末尾的3个frame [slog.XXX, slog.log, Logger.Handle]
			stacktrace.Frames = stacktrace.Frames[:ln-3]
		}
		record.AddAttrs(slog.Any(StacktraceKey, stacktrace))
	}

	if l.consoleLog.Enabled(ctx, record.Level) {
		if err := l.consoleLog.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}

	if l.sentryLog.Enabled(ctx, record.Level) {
		if err := l.sentryLog.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// WithAttrs implements slog.Handler
func (l *Logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Logger{
		consoleLog: l.consoleLog.WithAttrs(attrs),
		sentryLog:  l.sentryLog.WithAttrs(attrs),
	}
}

// WithGroup implements slog.Handler
func (l *Logger) WithGroup(name string) slog.Handler {
	return &Logger{
		consoleLog: l.consoleLog.WithGroup(name),
		sentryLog:  l.sentryLog.WithGroup(name),
	}
}

type Option struct {
	MinLogLevel     slog.Leveler
	StacktraceLevel slog.Leveler

	LogOption       slog.HandlerOptions
	SentryLogOption sentryslog.Option

	UseJson bool
}

func NewLogger(opt Option) *slog.Logger {
	var consoleLogger slog.Handler
	opt.LogOption.Level = opt.MinLogLevel
	opt.SentryLogOption.Level = slog.LevelError

	if opt.SentryLogOption.Converter == nil {
		opt.SentryLogOption.Converter = sentryslog.DefaultConverter
	}
	oldConvert := opt.SentryLogOption.Converter

	opt.SentryLogOption.Converter = func(
		addSource bool,
		replaceAttr func(groups []string, a slog.Attr) slog.Attr,
		loggerAttr []slog.Attr,
		groups []string,
		record *slog.Record,
		hub *sentry.Hub) *sentry.Event {

		var stacktrace *sentry.Stacktrace
		record.Attrs(func(a slog.Attr) bool {
			if a.Key == StacktraceKey {
				v, ok := a.Value.Any().(*sentry.Stacktrace)
				if ok {
					stacktrace = v
				}
				return false
			}
			return true
		})
		if stacktrace == nil {
			stacktrace = sentry.NewStacktrace()
		}

		event := oldConvert(addSource, replaceAttr, loggerAttr, groups, record, hub)
		event.Threads = append(event.Threads, sentry.Thread{
			Stacktrace: stacktrace,
		})
		return event
	}

	oldSentryReplaceAttr := opt.SentryLogOption.ReplaceAttr
	opt.SentryLogOption.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == StacktraceKey {
			return slog.Attr{}
		}
		if oldSentryReplaceAttr != nil {
			return oldSentryReplaceAttr(groups, a)
		}
		return a
	}

	oldReplaceAttr := opt.LogOption.ReplaceAttr
	opt.LogOption.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == StacktraceKey {
			val, _ := a.Value.Any().(*sentry.Stacktrace)
			if val == nil {
				// 不符合预期的情况, 将这个属性丢弃
				return slog.Attr{}
			}
			var sb strings.Builder
			for i := len(val.Frames) - 1; i >= 0; i-- {
				frame := val.Frames[i]
				sb.WriteString(frame.Module)
				sb.WriteByte('.')
				sb.WriteString(frame.Function)
				sb.WriteString("\n\t")
				sb.WriteString(findValieString("unknown", frame.Filename, frame.AbsPath))
				sb.WriteString(":")
				sb.WriteString(strconv.Itoa(frame.Lineno))
				if i != 0 {
					sb.WriteString("\n")
				}
			}
			return slog.Attr{
				Key:   a.Key,
				Value: slog.StringValue(sb.String()),
			}
		}
		if oldReplaceAttr != nil {
			return oldReplaceAttr(groups, a)
		}
		return a
	}

	if opt.UseJson {
		consoleLogger = slog.NewJSONHandler(os.Stdout, &opt.LogOption)
	} else {
		consoleLogger = slog.NewTextHandler(os.Stdout, &opt.LogOption)
	}
	return slog.New(&Logger{
		consoleLog:      consoleLogger,
		stacktraceLevel: opt.StacktraceLevel,
		sentryLog:       opt.SentryLogOption.NewSentryHandler(),
	})
}

func findValieString(defaultVal string, values ...string) string {
	for _, val := range values {
		if val != "" {
			return val
		}
	}
	return defaultVal
}
