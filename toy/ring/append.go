package ring

func CopyAppend(arr []int, v int) []int {
	copy(arr, arr[1:])
	arr[len(arr)-1] = v
	return arr
}

func ReSliceAppend(arr []int, v int) []int {
	arr[0] = 0
	arr = arr[1:]
	arr = append(arr, v)
	return arr
}

func Fill[T any](n int, v T) []T {
	var arr = make([]T, n)
	for i := range arr {
		arr[i] = v
	}
	return arr
}
