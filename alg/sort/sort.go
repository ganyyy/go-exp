package sort

/*
	是否稳定:
	最好时间复杂度:,
	平均时间复杂度:,
	最坏时间复杂度:
	空间复杂度:
*/

func BubbleSort(arr []int) {

	/*
		是否稳定:是
		最好时间复杂度:O(n^2),
		平均时间复杂度:O(n^2),
		最坏时间复杂度:O(n)
		空间复杂度:O(1)
	*/

	// 比较相邻的两个元素, 按照大小进行交换
	for i := 0; i < len(arr)-1; i++ { // 对比的次数
		// change 表示本次排序是否发生了变化
		var change bool
		// 将最大值放置到末尾
		for j := 0; j < len(arr)-i-1; j++ {
			if arr[j+1] >= arr[j] {
				continue
			}
			change = true
			arr[j+1], arr[j] = arr[j], arr[j+1]
		}
		if !change {
			break
		}
	}
}

func SelectSort(arr []int) {
	/*
		是否稳定: 否
		最好时间复杂度: O(n^2)
		平均时间复杂度: O(n^2)
		最坏时间复杂度: O(n^2)
		空间复杂度:  O(1)
	*/
	// 不稳定的原因: [5,8,5,2,6] -> [2,8,5,5,6].
	// 将整体分为两个区间, 已排序的区间 [:i], 未排序的区间[i:]
	// 每次都从未排序区间中选取一个最小值放入到已排序的区间内
	for i := 0; i < len(arr)-1; i++ {
		var minIdx = i
		// 筛选未排序区间的最小值
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		if i != minIdx {
			arr[i], arr[minIdx] = arr[minIdx], arr[i]
		}
	}
}

func InsertSort(arr []int) {
	/*
		是否稳定: 稳定(其实稳不稳定全看带不带等号)
		最好时间复杂度: O(n),
		平均时间复杂度: O(n^2),
		最坏时间复杂度: O(n^2)
		空间复杂度: O(1)
	*/
	// 整体分为两部分:
	// 已排序的部分, 未排序的部分.
	// 每次将未排序的部分中的首个元素插入到已排序部分中的指定位置上

	// 从1开始, 因为[:1]就是已排序的部分
	for i := 1; i < len(arr); i++ {
		var num = arr[i]
		// 从后向前迭代, 如果从前向后需要记录一个临时值
		var j = i - 1
		for ; j >= 0 && arr[j] > num; j-- {
			// 大于的直接交换位置
			arr[j+1] = arr[j]
		}
		// 将num插入到指定的位置上
		arr[j+1] = num
	}
}

func QuickSort(arr []int) {

	/*
		是否稳定: 不稳定. 最后一波交换位置可能会互换
		最好时间复杂度: O(n log(n)),
		平均时间复杂度: O(n log(n)),
		最坏时间复杂度: O(n^2)
		空间复杂度: O(log(n))~O(n)
	*/

	var sort func(int, int)

	sort = func(left, right int) {
		if left >= right {
			return
		}
		// 基准
		var pivot = arr[left]

		// 通过双指针 确定分区分界点
		var idx = left
		for i := left + 1; i <= right; i++ {
			if arr[i] < pivot {
				idx++
				arr[i], arr[idx] = arr[idx], arr[i]
			}
		}

		// 将基准点放置到合适的位置上
		arr[left], arr[idx] = arr[idx], arr[left]

		// 针对左半部分和右半部分进行排序
		sort(left, idx-1)
		sort(idx+1, right)
	}

	sort(0, len(arr)-1)
}

func ShellSort(arr []int) {
	/*
		是否稳定: 不稳定
		最好时间复杂度: O(n log^2(n)),
		平均时间复杂度: O(n log(n)),
		最坏时间复杂度: O(n log^2(n))
		空间复杂度: O(1)
	*/

	// 插入排序的优化
	// 选择增量gap, 分组插入排序
	var length = len(arr)
	var gap = length / 2

	for gap > 0 {
		for i := gap; i < length; i++ {
			var val = arr[i]
			var idx = i - gap
			// 插入排序
			for ; idx >= 0 && arr[idx] > val; idx = idx - gap {
				arr[idx+gap] = arr[idx]
			}
			arr[idx+gap] = val
		}
		gap /= 2
	}
}

func MergeSort(arr []int) {
	/*
		是否稳定: 稳定
		最好时间复杂度: O(n log(n)),
		平均时间复杂度: O(n log(n)),
		最坏时间复杂度: O(n log(n))
		空间复杂度: O(n)
	*/

	// 临时空间
	var space = make([]int, len(arr))

	// 将两个有序数组进行合并
	var merge = func(left, right []int) {
		// left的数据天然在顺序上早于right
		// 所以可以保证始终相对有序
		var li, ri int
		var idx int
		for li < len(left) && ri < len(right) {
			if left[li] <= right[ri] {
				space[idx] = left[li]
				li++
			} else {
				// 这里就可以记录一些骚操作.
				// 比如: 逆序对的个数!
				space[idx] = right[ri]
				ri++
			}
			idx++
		}
		for li < len(left) {
			space[idx] = left[li]
			li++
			idx++
		}
		for ri < len(right) {
			space[idx] = right[ri]
			ri++
			idx++
		}
	}

	var mergeSort func(int, int)

	mergeSort = func(start, end int) {
		if end-start < 2 {
			return
		}
		var middle = start + (end-start)/2
		mergeSort(start, middle)
		mergeSort(middle, end)
		merge(arr[start:middle], arr[middle:end])
		// 将排序后的数据更新到原始数组中
		copy(arr[start:end], space[:end-start])
		// log.Printf("arr:%v, space:%v", arr, space)
	}

	mergeSort(0, len(arr))
}

func HeapSort(arr []int) {

	var length = len(arr)
	if length <= 1 {
		return
	}

	var heap = func(idx int) {
		var parent = idx
		for {
			var child = parent*2 + 1
			if child >= length || child < 0 {
				break
			}
			// 选取较小的子节点
			if child+1 < length && arr[child+1] < arr[child] {
				child = child + 1
			}
			// 对比父节点和子节点的大小
			if arr[child] >= arr[parent] {
				break
			}
			// 交换父节点和较小的根节点之间的位置
			// 保证满足堆的性质
			arr[parent], arr[child] = arr[child], arr[parent]
			parent = child
		}
	}

	// 从中间节点开始, 构建小顶堆
	// 将整个数组看成是一个完全二叉树
	for i := length/2 - 1; i >= 0; i-- {
		heap(i)
	}

	// 堆头就是最小值, 然后递归找出接下来的最小值
	// 这是尾递归, 所以基本上不需要额外的栈空间
	HeapSort(arr[1:])
}
