// +build ignore

package main

import "fmt"

func main() {
	var src = []int{1, 2, 3, 4, 5, 6, 7}
	rotate(src, 4)
	fmt.Println(src)

	var ch = make(chan int)
	fmt.Println(len(ch))
}

func rotate(nums []int, k int) {
	nums = append(nums[k:], nums[:k]...)
	fmt.Println(nums)
}
