package generic2

import (
	"strconv"
)

//MyGeneric 带有底层类型约束的接口叫做 类型集, 只能用作泛型约束
type MyGeneric interface {
	~int | ~string  // 明确指定底层类型必须要是 int 或者 string
	String() string // 内嵌接口/函数都可以. 但是因为 存在了 类型约束, 所以此时该接口只能用来作为泛型约束参数
}

type Plus interface {
	Add(int) int
}

type MyNumber int

func (m MyNumber) String() string {
	return strconv.Itoa(int(m))
}

func (m MyNumber) Add(v int) int {
	return int(m) + v
}

type MyString string

func (m MyString) String() string {
	return string(m)
}

func RangeConstraints[T MyGeneric](src []T) {
	Range(src)  // 因为 泛型参数是 []MyGeneric, 所以其本身也实现了 fmt.Stringer 接口, 约束放宽是可以的
	Range2(src) // 任意类型皆可
}

//RangeConstraints1 src1, src2的类型必须要相同
func RangeConstraints1[T1 MyGeneric](src1, src2 []T1) {
	Range(src1)
	Range2(src2)
}

//RangeConstraints2 此时, 传入的两个参数可以拥有不同的类型
func RangeConstraints2[T1, T2 MyGeneric](src1 []T1, src2 []T2) {
	Range(src1)
	Range2(src2)
}

//RangeConstraints3 允许声明多种类型参数
func RangeConstraints3[T1 MyGeneric, T2 Plus](src1 []T1, src2 []T2) {
	Range(src1)
	var ret int
	for _, v := range src2 {
		ret = v.Add(ret)
	}
}
