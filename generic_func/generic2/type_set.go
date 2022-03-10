package generic2

import "golang.org/x/exp/constraints"

/*
	1. Go不支持运算符重载
	2. 提供了内置的 comparable 泛型类型来约束map的key的类型
	3. 提供了内置的 constraints.XXX, 用来约束可以进行比较的类型
*/

func Small[T constraints.Ordered](src []T) T {
	if len(src) == 0 {
		var ret T
		return ret
	}
	var min = src[0]

	for _, v := range src[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

/*
	1. 普通的 interface 是一个无限集, 只要满足了其中的方法就满足约束. 就约束而言, 这个集合只有它一个元素
	2. 类型集约束相当于在 原有的 interface 上, 添加了底层类型约束/方法约束
*/

type M1 interface {
	MM1()
}

type M2 interface {
	MM2()
}

//MMSet 要求必须要同时实现 MM1和MM2两个方法
type MMSet interface {
	M1
	M2
}

//MMSetWithMethod 还要求实现Show方法
type MMSetWithMethod interface {
	MMSet
	Show()
}

var (
	_ M1
	_ M2
	_ MMSet
	_ MMSetWithMethod
	//_ MMSetMethodWithType // 只能用在约束中
)

// 以上都是原有的interface

//MMSetMethodWithType 再增加了底层类型约束, 此时就不能当成传统的接口使用了
type MMSetMethodWithType interface {
	~int // 底层类型约束只能是 原始类型(包括map[K]V, []T, chan等内置类型). 不允许使用自定义的类型作为底层类型约束
	MMSetWithMethod
}

type MMUnionSet interface {
	~string | ~float64 // 底层类型可以是 string 或者 float64
	MMSetWithMethod
}

//Index 可比较的类型约束, 仅用于 ==和!=
func Index[T comparable](s []T, v T) int {
	if len(s) == 0 {
		return -1
	}
	for i, vv := range s {
		if vv == v {
			return i
		}
	}
	return -1
}

type ImpossibleConstraint interface {
	comparable
	[]int // 切片类型不可比较, 所以不存在 满足该约束的类型
}

type Comparable interface {
	comparable
}

func Equal[T Comparable](a, b T) bool {
	return a == b
}
