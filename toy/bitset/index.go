package bitset

import (
	"math/bits"
)

type Index uint32

func (i Index) Row(shift Shift) uint {
	return (uint(i) >> shift) & ((1 << Bits) - 1)
}

func (i Index) Offset(shift Shift) uint {
	return (uint(i)) / (1 << shift)
}

func (i Index) Mask(shift Shift) uint {
	return 1 << i.Row(shift)
}

func Offsets(bit Index) (off1, off2, off3 uint) {
	return bit.Offset(Shift1), bit.Offset(Shift2), bit.Offset(Shift3)
}

func AverageOnes(v uint) uint {
	if v <= 1 {
		return 0
	}
	return uint(bits.Len(v)) - 1
}
