package timewheel_test

import (
	"testing"
	"unsafe"

	"github.com/go-playground/assert/v2"
)

var srcSource = []string{
	"1",
	"2",
}

var group = [][]string{
	{"1", "2"},
	{"1", "1", "2", "2"},
	{"1", "1", "1", "2", "2", "2"},
}

func Test123(t *testing.T) {
	var src uint64
	var cnt [8]uint8

	assert.Equal(t, src, *(*uint64)(unsafe.Pointer(&cnt[0])))
}
