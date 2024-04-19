package common

import "context"

var (
	runnerIdxKey struct{}
)

func WithRunnerIdx(ctx context.Context, idx int) context.Context {
	return context.WithValue(ctx, runnerIdxKey, idx)
}

func RunnerIdxFromContext(ctx context.Context) (int, bool) {
	idx, ok := ctx.Value(runnerIdxKey).(int)
	return idx, ok
}
