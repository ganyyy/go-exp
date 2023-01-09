package api

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func ArrayToString[T any](arr []T) (ret string) {
	ln := len(arr)
	if ln == 0 {
		return
	}
	var ele T
	size := int(unsafe.Sizeof(ele))
	if size == 0 {
		return
	}
	var data = unsafe.SliceData(arr)
	return unsafe.String((*byte)(unsafe.Pointer(data)), len(arr)*size)
}

func ArrayEqual[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	return ArrayToString(a) == ArrayToString(b)
}

func TestArrayEqual(t *testing.T) {
	genSlice := func(num int) (ret []int) {
		for i := 0; i < num; i++ {
			ret = append(ret, 0)
		}
		return ret
	}

	a, b := genSlice(10), genSlice(10)
	assert.True(t, ArrayEqual(a, b))
}
