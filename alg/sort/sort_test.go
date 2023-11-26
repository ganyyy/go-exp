package sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

type SortFunc func([]int)

func shuffle(arr []int) []int {
	var old = make([]int, len(arr))
	r.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	copy(old, arr)
	return old
}

const NUM = 1000

func checkSort(t *testing.T, fn SortFunc) {
	var arr = make([]int, NUM)
	for i := range arr {
		arr[i] = i
	}
	var old = shuffle(arr)
	var backup = make([]int, len(old))
	copy(backup, old)

	fn(arr)
	sort.Ints(old)
	assert.Equal(t, old, arr, "src:%v", backup)
}

func Test_SortFunc(t *testing.T) {
	tests := []struct {
		name string
		fn   SortFunc
	}{
		{"BubbleSort", BubbleSort},
		{"SelectSort", SelectSort},
		{"InsertSort", InsertSort},
		{"QuickSort", QuickSort},
		{"ShellSort", ShellSort},
		{"MergeSort", MergeSort},
		{"HeapSort", HeapSort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkSort(t, tt.fn)
		})
	}
}

func BenchmarkSortFunc(b *testing.B) {
	tests := []struct {
		name string
		fn   SortFunc
	}{
		{"BubbleSort", BubbleSort},
		{"SelectSort", SelectSort},
		{"InsertSort", InsertSort},
		{"QuickSort", QuickSort},
		{"ShellSort", ShellSort},
		{"MergeSort", MergeSort},
		{"HeapSort", HeapSort},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			var arr = make([]int, NUM)
			for i := range arr {
				arr[i] = i
			}
			var old = shuffle(arr)
			for i := 0; i < b.N; i++ {
				copy(arr, old)
				tt.fn(arr)
			}
		})
	}
}
