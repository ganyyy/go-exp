package bitset

import (
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBits(t *testing.T) {
	t.Run("Leading", func(t *testing.T) {
		var vals = []uint{
			0b0,
			0b1,
			0b10,
			0b11,
			0b101,
			0b110,
		}

		for _, v := range vals {
			t.Logf("=========%b=========", v)
			t.Logf("%v", bits.LeadingZeros(v))
			t.Logf("%v", bits.TrailingZeros(v))
			t.Logf("%v", bits.Len(v))
		}
	})
}

func TestAverageOnes(t *testing.T) {
	var vals = []uint{
		0b0,
		0b1,
		0b10,
		0b101, 0b10110,
		0b100010,
	}

	for _, v := range vals {
		t.Logf("%v", AverageOnes(v))
	}
}

func TestIndex(t *testing.T) {

	var shifts = []Shift{
		Shift0,
		Shift1,
		Shift2,
		Shift3,
	}

	var vals = [...]uint{
		1 << Shift1,
		1<<Shift1 + 1,
		1<<Shift2 + 1,
		1 << Shift3,
	}

	t.Log(vals)

	for _, v := range vals {
		idx := Index(v)
		t.Log(Offsets(idx))

		for _, shift := range shifts {
			t.Log("shift:", idx, idx.Mask(shift), idx.Offset(shift), idx.Row(shift))
		}
	}
}

var genId = func() uint {
	return uint(rand.Uint64() % (MaxEID))
}

func TestBitSet(t *testing.T) {

	var bs BitSet
	for i := 0; i < 10000; i++ {
		idx := Index(genId())
		assert.Equal(t, bs.Contain(idx), bs.Add(idx))
		assert.True(t, bs.Add(idx))
		assert.True(t, bs.Remove(idx))
	}
}

func BenchmarkBitSet(b *testing.B) {
	var bs BitSet
	for i := 0; i < b.N; i++ {
		idx := Index(genId())
		bs.Add(idx)
	}
}

func TestRange(t *testing.T) {
	var bs BitSet
	for i := 0; i < 10; i++ {
		bs.Add(Index(i * 100))
		bs.Add(Index(MaxEID - i*100 - 1))
	}

	for idx := range bs.Range() {
		t.Logf("%v", idx)
	}
}
