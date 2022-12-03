package fuzzing

import (
	"math"
)

func ifCount(n int) (ret int) {
	if n == 0 {
		return
	}
	// 先转成float, 否则会有溢出的情况
	v := float64(n)
	if v < 0 {
		ret++
		v = -v
	}

	if v < 1e1 {
		ret += 1
	} else if v < 1e2 {
		ret += 2
	} else if v < 1e3 {
		ret += 3
	} else if v < 1e4 {
		ret += 4
	} else if v < 1e5 {
		ret += 5
	} else if v < 1e6 {
		ret += 6
	} else if v < 1e7 {
		ret += 7
	} else if v < 1e8 {
		ret += 8
	} else if v < 1e9 {
		ret += 9
	} else if v < 1e10 {
		ret += 10
	} else if v < 1e11 {
		ret += 11
	} else if v < 1e12 {
		ret += 12
	} else if v < 1e13 {
		ret += 13
	} else if v < 1e14 {
		ret += 14
	} else if v < 1e15 {
		ret += 15
	} else if v < 1e16 {
		ret += 16
	} else if v < 1e17 {
		ret += 17
	} else if v < 1e18 {
		ret += 18
	} else {
		ret += 19
	}
	return
}

func loopCount(v int) int {
	var n int
	if v < 0 {
		n++
	}
	for v != 0 {
		v /= 10
		n++
	}
	return n
}

func logCount(v int) int {
	if v == 0 {
		return 0
	}
	var n int
	// 先转成float, 否则会出现溢出的情况
	f := float64(v)
	if f < 0 {
		n++
		f = -f
	}
	n += int(math.Log10(float64(f))) + 1
	return n
}
