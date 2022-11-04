package generic2

import (
	"strconv"
	"testing"
)

func TestMapReduce(t *testing.T) {
	var src = []string{
		"1", "2", "3", "4", "5", "6", "7",
	}

	// 简单的函数调用可以自动推断类型

	// 自动推断为 Map[string, string]
	var _ = Map(src, func(f string) string {
		return f
	})

	// 不知道咋回事, 太过深层次的函数调用无法自动推断
	// 如果无法推断出类型, 那么就会编译失败
	// 现在IDE的自动推断问题太大了, 我的评价是: 不如Rust
	var ret = Reduce(
		Filter(
			Map(
				src,
				func(v string) int {
					var ret, _ = strconv.Atoi(v)
					return ret
				}),
			func(v int) bool {
				return v&1 == 0
			}),
		func(a, b int) int {
			return a + b
		}, 100)

	t.Logf("%+v", ret)
}

func TestShowTypeInterface(t *testing.T) {
	var src []map[int]bool
	ShowTypeInterface1(src) // 原类型匹配
	ShowTypeInterface2(src) // 泛型匹配
	ShowTypeInterface3(src) // T1 = map[int]bool
	ShowTypeInterface4(src) // T1 = int, T2 = bool
}
