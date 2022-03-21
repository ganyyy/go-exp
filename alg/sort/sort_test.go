package sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SortFunc func([]int)

func shuffle(arr []int) []int {
	var old = make([]int, len(arr))
	rand.Shuffle(len(arr), func(i, j int) {
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

func TestBubbleSort(t *testing.T) {
	checkSort(t, BubbleSort)
}

func TestSelectSort(t *testing.T) {
	checkSort(t, SelectSort)
}

func TestInsertSort(t *testing.T) {
	checkSort(t, InsertSort)
}

func TestQuickSort(t *testing.T) {
	checkSort(t, QuickSort)
}

func TestShellSort(t *testing.T) {
	checkSort(t, ShellSort)
}

func TestMergeSort(t *testing.T) {
	checkSort(t, MergeSort)
}

func TestHeapSort(t *testing.T) {
	checkSort(t, HeapSort)
}
