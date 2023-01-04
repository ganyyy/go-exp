package bitset

import "fmt"

const (
	Bits = 6 // 每一层使用的Bit位. 因为2^6 == 64, 不能再大了

	Layers = 4 // 层数(数值分层). 层次可以累加, 但是逻辑也需要大改

	Max = Bits * Layers // 最长的Bit位

	MaxEID = 1 << Max // 位图最多支持的元素的个数

	DefaultVal = 0
)

type Shift uint

const (
	Shift0 Shift = iota * Bits
	Shift1
	Shift2
	Shift3
	ShiftCount Shift = iota
)

// 一些常量检查函数

func checkShift(shift Shift) {
	if shift >= ShiftCount {
		panic(fmt.Sprintf("invalid shift %v", shift))
	}
}

func validRange(idx Index) {
	if idx >= MaxEID {
		panic(fmt.Sprintf("MAX_EID is %v, found %v", MaxEID-1, idx))
	}
}
