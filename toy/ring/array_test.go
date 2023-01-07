package ring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayRing(t *testing.T) {
	const SIZE = 10
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

	for i := 0; i < 2; i++ {
		buffer.Clear()
		for i := 0; i < SIZE*5; i++ {
			buffer.Add(i)
			t.Logf("%v: %+v", i, buffer)
			cp := buffer.Copy()
			ln := buffer.Len()
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
}

func BenchmarkRingAdd(b *testing.B) {
	var buffer = NewArrayRing[int](10)
	for i := 0; i < b.N; i++ {
		buffer.Add(i)
	}
}
