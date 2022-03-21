package sync_map_test

import (
	"math"
	"testing"

	"ganyyy.com/go-exp/container/sync_map"
	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	var m = sync_map.NewSyncMap[int, int]()

	m.Store(100, 200)
	var v, ok = m.Load(100)
	assert.Equal(t, 200, v)
	assert.True(t, ok)

	v, ok = m.Load(200)
	assert.Equal(t, v, 0)
	assert.False(t, ok)

	v, ok = m.LoadAndDelete(100)
	assert.Equal(t, 200, v)
	assert.True(t, ok)

	v, ok = m.LoadAndDelete(200)
	assert.Equal(t, 0, v)
	assert.False(t, ok)

	m.Store(100, 200)
	v, ok = m.LoadOrStore(100, 300)
	assert.Equal(t, v, 200)
	assert.True(t, ok)

	v, ok = m.LoadOrStore(200, 300)
	assert.Equal(t, v, 300)
	assert.False(t, ok)

	m.Delete(200)
	v, ok = m.Load(200)
	assert.Equal(t, v, 0)
	assert.False(t, ok)

	m.Range(func(i1, i2 int) bool {
		t.Logf("key %v, value %v", i1, i2)
		m.Delete(i1)
		return true
	})

	var check = make(map[int]bool)
	for i := 0; i < 100; i++ {
		m.Store(i, i+100)
		check[i] = true
	}
	m.Range(func(i1, i2 int) bool {
		delete(check, i1)
		return true
	})
	assert.Equal(t, 0, len(check))

}

func TestNaNMap(t *testing.T) {
	var m = sync_map.NewSyncMap[float64, string]()

	m.Store(math.NaN(), "100")
	m.Store(math.NaN(), "200")
	m.Store(math.NaN(), "300")
	m.Store(math.NaN(), "400")
	m.Store(math.NaN(), "500")
	m.Store(math.NaN(), "600")
	m.Store(math.NaN(), "700")
	m.Store(math.NaN(), "800")

	m.Delete(math.NaN())

	m.Range(func(f float64, s string) bool {
		t.Logf("%v, %v, %v", f, s, math.IsNaN(f))
		return true
	})

}
