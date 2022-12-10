package api

import (
	"arena"
	"errors"
	"runtime"
	"testing"
	"time"
)

func TestErrorWrap(t *testing.T) {
	err := errors.Join(
		errors.New("1"),
		errors.New("1"),
		errors.New("1"),
		errors.New("1"),
		errors.New("1"),
	)

	t.Logf("err:%+v", err)

	errors.Unwrap(err)

	errors.Is(err, nil)
}

func TestTimeFormat(t *testing.T) {
	var nowt = time.Now()

	t.Logf("%+v", nowt.Format(time.DateOnly))
}

func TestArena(t *testing.T) {
	a := arena.NewArena()
	defer a.Free()
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	for i := 0; i < 2<<10; i++ {
		_ = arena.New[int](a)
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	arena.MakeSlice[int](a, 100, 200)

	t.Logf("total alloc:%v", after.TotalAlloc-memStat.TotalAlloc)
	// t.Logf("before:%+v", memStat)
	// t.Logf("after:%+v", after)
}

func Error[T error](es ...T) string {
	var bs []byte
	for i, e := range es {
		if i > 0 {
			bs = append(bs, '\n')
		}
		bs = append(bs, e.Error()...)
	}
	return string(bs)
}
