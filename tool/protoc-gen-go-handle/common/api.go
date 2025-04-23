package common

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
)

type ClientInterface interface {
	Invoke(context.Context, string, proto.Message, proto.Message) error
}

type InvokeOpt struct {
	ctx     context.Context
	timeout time.Duration
}

type ApplyOption func(*InvokeOpt)

func WithContext(ctx context.Context) ApplyOption {
	return func(opt *InvokeOpt) {
		opt.ctx = ctx
	}
}

func WithTimeout(d time.Duration) ApplyOption {
	return func(opt *InvokeOpt) {
		opt.timeout = d
	}
}

func GenerateOptions(opts ...ApplyOption) InvokeOpt {
	opt := InvokeOpt{
		timeout: time.Second,
		ctx:     context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func (opt InvokeOpt) Context() (context.Context, context.CancelFunc) {
	ctx := opt.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, opt.timeout)
}
