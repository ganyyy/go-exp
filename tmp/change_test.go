package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNilSlice(t *testing.T) {
	var s1 = make([]int, 0)
	var s2 []int
	var s3 = *(*[]int)(unsafe.Pointer(&reflect.SliceHeader{}))

	assert.Equal(t, s2, s3)
	assert.NotEqual(t, s1, s2)
	assert.Equal(t, []int(nil), s2)
	assert.NotEqual(t, []int(nil), s1)

	println(s1)
	println(s2)
	println(s3)
}
