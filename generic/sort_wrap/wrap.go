package sortwrap

import "sort"

type sortWrapper[T any] struct {
	arr  []T
	less func(T, T) bool
}

func (s *sortWrapper[T]) Len() int {
	return len(s.arr)
}

func (s *sortWrapper[T]) Swap(i, j int) {
	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
}

func (s *sortWrapper[T]) Less(i, j int) bool {
	return s.less(s.arr[i], s.arr[j])
}

func Sort[T any](arr []T, cmp func(T, T) bool) {
	sort.Sort(&sortWrapper[T]{
		arr:  arr,
		less: cmp,
	})
}
