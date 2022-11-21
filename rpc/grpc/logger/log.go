package logger

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
)

type logger struct {
	*log.Logger
}

func newLogger() io.Writer {
	return &logger{
		Logger: log.New(os.Stdout, "[GRPC]", 0),
	}
}

func (l *logger) Write(data []byte) (int, error) {
	l.Printf("%s", data)
	return len(data), nil
}

func SetGRPCLogger() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(newLogger(), newLogger(), newLogger()))
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
	log.Printf("[%v] HandleRPC info:%+v", h.tag, info)
}
func (h *handle) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	log.Printf("[%v] TagConn info:%+v", h.tag, info)
	return ctx
}
func (h *handle) HandleConn(ctx context.Context, info stats.ConnStats) {
	log.Printf("[%v] HandleConn info:%+v", h.tag, info)
}
