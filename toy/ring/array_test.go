package ring

import (
	"math"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAppend(b *testing.B) {

	const S = 1000

	gen := func(SIZE int) []int {
		var arr = make([]int, 0, SIZE)
		for i := 0; i < SIZE; i++ {
			arr = append(arr, i)
		}
		return arr
	}
	getHeader := func(arr []int) uintptr {
		header := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
		return header.Data
	}
	_ = getHeader
	b.Run("in place", func(b *testing.B) {
		arr := gen(S)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			arr = CopyAppend(arr, i)
		}
		// b.Logf("%v, %X", cap(arr), getHeader(arr))
	})
	b.Run("append", func(b *testing.B) {
		arr := gen(S)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			arr = ReSliceAppend(arr, i)
		}
		// b.Logf("%v, %X", cap(arr), getHeader(arr))
	})
	b.Run("Copy", func(b *testing.B) {
		const CAP = S
		arr := gen(CAP)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			copy(arr, arr[1:])
			arr[CAP-1] = i
		}
	})
	b.Run("Ring", func(b *testing.B) {
		var buffer = NewArrayRing[int](S)
		for i := 0; i < b.N; i++ {
			buffer.Add(i)
		}
	})
}

func TestArrayRing(t *testing.T) {
	const SIZE = 10
	t.Run("normal", func(t *testing.T) {
		var buffer = NewArrayRing[int](SIZE)

		func() {
			defer func() {
				assert.NotNil(t, recover())
			}()
			NewArrayRing[int](-1)
		}()

		{
			_, ok := buffer.Get(0)
			assert.False(t, ok)
			_, ok = buffer.Get(-1)
			assert.False(t, ok)
		}

		for i := 0; i < 5; i++ {
			buffer.Clear()
			for i := 0; i < SIZE*4; i++ {
				buffer.Add(i)
				cp := buffer.Copy()
				ln := buffer.Len()
				t.Logf("%v: %+v, cp:%v, ln:%v", i, buffer, cp, ln)
				assert.Equal(t, ln, len(cp))
				for i, v := range cp {
					getV, ok := buffer.Get(i)
					getNV, nvOk := buffer.Get(-ln + i)
					assert.True(t, ok)
					assert.True(t, nvOk)
					assert.Equal(t, getV, v)
					assert.Equal(t, getV, getNV)
				}
			}
		}

		{

			_, ok := buffer.Get(SIZE)
			assert.False(t, ok)
			_, ok = buffer.Get(-SIZE - 1)
			assert.False(t, ok)
		}

		{
			buffer.Add(100)
			buffer.Range(nil)
			var rangeV int
			buffer.Range(func(i int) bool {
				rangeV = i
				return false
			})
			getV, ok := buffer.Get(0)
			assert.Equal(t, rangeV, getV)
			assert.True(t, ok)
		}

		{
			buffer.Clear()
			assert.Equal(t, buffer.Len(), 0)
		}

		{
			_, ok := buffer.Get(math.MinInt64)
			assert.False(t, ok)
		}
	})

	t.Run("OneElement", func(t *testing.T) {
		var buffer = NewArrayRing[int](SIZE)
		buffer.Add(1)
		t.Logf("%+v, ln:%v, cp:%v", buffer, buffer.Len(), buffer.Copy())
	})

	t.Run("PopTop", func(t *testing.T) {
		var buffer = NewArrayRing[int](SIZE)

		check := func(wantV int, wantOk bool) {
			tv, tOk := buffer.Top()
			pv, pOk := buffer.Pop()
			assert.Equal(t, tv, pv)
			assert.Equal(t, tOk, pOk)
			assert.Equal(t, tv, wantV)
			assert.Equal(t, tOk, wantOk)
		}
		check(0, false)

		for i := 0; i < SIZE*2; i++ {
			buffer.Add(i)
			check(i, true)
		}

		buffer.Clear()

		for i := 1; i <= SIZE*2; i++ {
			buffer.Add(i)

			if i%SIZE == 0 {
				for i := 0; i < 5 && buffer.Len() != 0; i++ {
					t.Logf("%+v, %v", buffer, buffer.Copy())
					check(buffer.Get(0))
					buffer.Add(i)
				}
				ln := buffer.Len()
				for buffer.Len() != 0 {
					check(buffer.Get(0))
					ln--
					assert.Equal(t, ln, buffer.Len())
				}
				t.Logf("%+v, %v", buffer, buffer.Copy())
			}
		}
		t.Logf("%+v, %v", buffer, buffer.Copy())
	})
}
