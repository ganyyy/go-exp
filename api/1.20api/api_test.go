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
	log := func(from string, before runtime.MemStats) runtime.MemStats {
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		t.Logf("[%s] total alloc:%v", from, int(after.TotalAlloc)-int(before.TotalAlloc))
		return after
	}
	var memStat = log("init", runtime.MemStats{})
	var tmp []*int
	var tmp2 []int
	for i := 0; i < 2<<10; i++ {
		tmp = append(tmp, arena.New[int](a))
	}
	memStat = log("new int", memStat)
	tmp2 = arena.MakeSlice[int](a, 100, 2000)
	memStat = log("new slice", memStat)
	runtime.GC()
	runtime.GC()
	memStat = log("alloc gc", memStat)
	a.Free()
	memStat = log("free", memStat)
	runtime.GC()
	runtime.GC()
	log("free gc", memStat)

	runtime.KeepAlive(tmp)
	runtime.KeepAlive(tmp2)
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
