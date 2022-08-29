package api_119

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"runtime/debug"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmtAppend(t *testing.T) {
	var buf [100]byte

	var show = fmt.Appendf(buf[:0], "%s", "hello world") // hello world
	show = fmt.Appendln(show, "1", "2")                  // hello world1 2
	show = fmt.Append(show, "3", "4")                    // hello world1 2\n3 4

	t.Logf("%s", show)
}

func TestMemoryLimit(t *testing.T) {
	debug.SetMemoryLimit(2 ^ 20) //GB以下不进行GC
}

type Common struct {
	Name string
}

func TestAtomic(t *testing.T) {
	var pc atomic.Pointer[Common]
	var ov = &Common{
		Name: "123",
	}
	pc.Store(ov)

	var v = pc.Load()

	var pi atomic.Uint64
	pi.Store(100)
	pi.Load()

	assert.Equal(t, v, ov)
}

func TodayHour(t time.Time, hour int) time.Time {
	var loc = time.Local
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, loc)
}

func TestAbs(t *testing.T) {
	var now = time.Now()
	var nextDay = now.AddDate(0, 0, 1)

	assert.Equal(t, now.Sub(nextDay).Abs(), nextDay.Sub(now).Abs())
	t.Logf("sub1:%v, sub2:%v", now.Sub(nextDay).Abs(), nextDay.Sub(now).Abs())
}

func TestTodayHour(t *testing.T) {
	var now = time.Now()

	var t1 = TodayHour(now, 5)
	var t2 = TodayHourMath(now, 5)
	assert.Equal(t, t1, t2)
	t.Logf("now:%v, now 5:%v", now, TodayHour(now, 5))
}

var zero = time.Date(2006, 1, 1, 0, 0, 0, 0, time.Local).Unix()

func TodayHourMath(t time.Time, hour int) time.Time {
	const HourSecond = 60 * 60
	const DaySecond = HourSecond * 24
	var subDay = (t.Unix() - zero) / DaySecond
	return time.Unix(zero+subDay*DaySecond+HourSecond*int64(hour), 0)
}

func BenchmarkLocal(b *testing.B) {
	var now = time.Now()
	b.Run("time local", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TodayHour(now, 5)
		}
	})

	b.Run("math local", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TodayHourMath(now, 5)
		}
	})
}
func TestMapHash(t *testing.T) {
	var seed = maphash.MakeSeed()

	var buf [16]byte
	_, _ = rand.Read(buf[:])
	for i := 0; i < 10; i++ {
		t.Logf("src:%s, hash:%v", buf, maphash.Bytes(seed, buf[:]))
	}
}
