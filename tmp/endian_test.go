package main

import (
	"sort"
	"strconv"
	"testing"
	"unsafe"
)

func TestEndian(t *testing.T) {
	var v int16 = 0x1234
	var arr = *(*[2]byte)(unsafe.Pointer(&v))
	t.Logf("%X, %X", arr[0], arr[1])
}

func TestNumberConvert(t *testing.T) {
	var v, _ = strconv.ParseInt(strconv.Itoa(100), 7, 64)

	t.Logf(strconv.Itoa(int(v)))
}

func TestCloseChannel(t *testing.T) {
	var taskChan = make(chan int, 3)

	close(taskChan)

	select {
	case taskChan <- 1:
	default:
	}
}

func TestSearchInts(t *testing.T) {
	var src = []int{1, 3, 5, 7, 9}
	t.Logf("%v", sort.SearchInts(src, 0))
	t.Logf("%v", sort.SearchInts(src, 2))
	t.Logf("%v", sort.SearchInts(src, 10)-1)
}
