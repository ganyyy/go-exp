package tower

import (
	"context"
	"fmt"
	"time"
)

type Timeout[Req, Resp any] struct {
	timeout time.Duration
	inner   S[Req, Resp]
}

func NewTimeout[Req, Resp any](timeout time.Duration) *Timeout[Req, Resp] {
	return &Timeout[Req, Resp]{timeout: timeout}
}

func (t *Timeout[Req, Resp]) Layer(svc S[Req, Resp]) S[Req, Resp] {
	t.inner = svc
	return t
}

func (t *Timeout[Req, Resp]) Call(ctx context.Context, req Req) (Resp, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()
	return t.inner.Call(ctx, req)
}

type Logger[Req, Resp any] struct {
	inner  S[Req, Resp]
	prefix string
}

func NewLogger[Req, Resp any](prefix string) *Logger[Req, Resp] {
	return &Logger[Req, Resp]{
		prefix: prefix,
	}
}

func (l *Logger[Req, Resp]) Layer(svc S[Req, Resp]) S[Req, Resp] {
	l.inner = svc
	return l
}

func (l *Logger[Req, Resp]) Call(ctx context.Context, req Req) (resp Resp, err error) {
	start := time.Now()

	fmt.Printf("[%s]%s Request received %+v\n", l.prefix, start.Format(time.RFC3339), req)
	defer func(start time.Time) {
		fmt.Printf("[%s]%s Response sent %+v, duration: %s\n", l.prefix, time.Now().Format(time.RFC3339), resp, time.Since(start))
	}(start)
	return l.inner.Call(ctx, req)
}

type Recovery[Req, Resp any] struct {
	inner S[Req, Resp]
}

func NewRecovery[Req, Resp any]() *Recovery[Req, Resp] {
	return &Recovery[Req, Resp]{}
}

func (r *Recovery[Req, Resp]) Layer(svc S[Req, Resp]) S[Req, Resp] {
	r.inner = svc
	return r
}

func (r *Recovery[Req, Resp]) Call(ctx context.Context, req Req) (resp Resp, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Recovered from panic: %v\n", e)
			err = fmt.Errorf("internal server error %+v", e)
		}
	}()

	return r.inner.Call(ctx, req)
}
