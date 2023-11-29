package logger

import (
	"context"
	"io"
	"log"
	"log/slog"

	"ganyyy.com/go-exp/helper"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
)

type logger struct {
	*slog.Logger
	level slog.Level
}

func newLogger(level slog.Level) io.Writer {
	return &logger{
		level:  level,
		Logger: helper.InitSlog(),
	}
}

func (l *logger) Write(data []byte) (int, error) {
	l.Log(context.Background(), l.level, string(data))
	return len(data), nil
}

func SetGRPCLogger() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		newLogger(slog.LevelInfo),
		newLogger(slog.LevelWarn),
		newLogger(slog.LevelError)))
}

type handle struct {
	tag string
}

func NewHandle(reason string) stats.Handler {
	return &handle{tag: reason}
}

func (h *handle) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	log.Printf("[%v] TagRPC info:%+v", h.tag, info)
	return ctx
}
func (h *handle) HandleRPC(ctx context.Context, info stats.RPCStats) {
	log.Printf("[%v] HandleRPC info [%T]:%+v", h.tag, info, info)
}
func (h *handle) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	log.Printf("[%v] TagConn info:%+v", h.tag, info)
	return ctx
}
func (h *handle) HandleConn(ctx context.Context, info stats.ConnStats) {
	log.Printf("[%v] HandleConn info:%+v", h.tag, info)
}
