package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"strings"

	"ganyyy.com/go-exp/helper"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
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

var putKey int
var rangeKey int
var watchKey int

type PutCtx struct {
	Key string
}

type RangeCtx struct {
	Key      string
	IsPrefix bool
}

type WatchKey struct {
	Key      string
	Revision int64
	IsPrefix bool
}

type WatchCtx struct {
	Keys map[int64]WatchKey
}

func (h *handle) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	// log.Printf("[%v] TagRPC info:%+v", h.tag, info)
	if strings.HasSuffix(info.FullMethodName, "Put") {
		return context.WithValue(ctx, &putKey, &PutCtx{})
	} else if strings.HasSuffix(info.FullMethodName, "Range") {
		return context.WithValue(ctx, &rangeKey, &RangeCtx{})
	} else if strings.HasSuffix(info.FullMethodName, "Watch") {
		return context.WithValue(ctx, &watchKey, &WatchCtx{Keys: make(map[int64]WatchKey)})
	}
	// return ctx
	return ctx
}
func (h *handle) HandleRPC(ctx context.Context, info stats.RPCStats) {
	// log.Printf("[%v] HandleRPC info [%T]:%+v", h.tag, info, info)

	putCtx, isPut := ctx.Value(&putKey).(*PutCtx)
	rangeCtx, isRange := ctx.Value(&rangeKey).(*RangeCtx)
	// watchCtx, isWatch := ctx.Value(&watchKey).(*WatchCtx)

	header := func() string {
		if isPut {
			return "Put"
		}
		if isRange {
			return "Get"
		}
		// if isWatch {
		// 	return "Watch"
		// }
		return "Unknown"
	}

	switch typ := info.(type) {
	case *stats.Begin:
	case *stats.PickerUpdated:
	case *stats.OutHeader:
	case *stats.OutPayload:
		if isRange {
			rangeRequest := typ.Payload.(*etcdserverpb.RangeRequest)
			rangeCtx.Key = string(rangeRequest.Key)
			rangeCtx.IsPrefix = rangeRequest.RangeEnd != nil
		} else if isPut {
			putRequest := typ.Payload.(*etcdserverpb.PutRequest)
			putCtx.Key = string(putRequest.Key)
			// } else if isWatch {
			// watchRequest := typ.Payload.(*etcdserverpb.WatchRequest)
			// createRequest := watchRequest.GetCreateRequest()
			// cancelRequest := watchRequest.GetCancelRequest()
			// if createRequest != nil {
			// 	var key WatchKey
			// 	key.Key = string(createRequest.Key)
			// 	key.IsPrefix = createRequest.RangeEnd != nil
			// 	key.Revision = createRequest.StartRevision
			// 	watchCtx.Keys[createRequest.WatchId] = key
			// 	log.Printf("Watch  create request:%+v", createRequest)
			// }
			// if cancelRequest != nil {
			// 	delete(watchCtx.Keys, cancelRequest.WatchId)
			// 	log.Printf("Watch cancel request:%+v", cancelRequest)
			// }
		}
		// log.Printf("[%v]OUT, payload %T %+v", header(), typ.Payload, typ.Payload)
	case *stats.InHeader:
	case *stats.InTrailer:
	case *stats.InPayload:
		// log.Printf("[%v]IN, payload %T %+v", header(), typ.Payload, typ.Payload)
	case *stats.End:
		// log.Printf("[%v]END, isClient %v, cost %v", header(), typ.IsClient(), typ.EndTime.Sub(typ.BeginTime))
		cost := typ.EndTime.Sub(typ.BeginTime)
		if isRange {
			log.Printf("[%v] key:%v, prefix:%v cost:%v", header(),
				rangeCtx.Key, rangeCtx.IsPrefix, cost)
		} else if isPut {
			log.Printf("[%v] key:%v cost:%v", header(),
				putCtx.Key, cost)
			// } else if isWatch {
			// 	log.Printf("[%v] stop cost:%v", header(), cost)
		}
	}

	// return
}
func (h *handle) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// log.Printf("[%v] TagConn info:%+v", h.tag, info)
	return ctx
}
func (h *handle) HandleConn(ctx context.Context, info stats.ConnStats) {
	// log.Printf("[%v] HandleConn info:%+v", h.tag, info)
}
