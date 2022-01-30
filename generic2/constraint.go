package generic2

import (
	"strconv"
)

type MyGeneric interface {
	//fmt.Stringer   // 明确指定要实现的方法
	~int | ~string  // 明确指定底层类型必须要是 int 或者 string
	String() string // 内嵌接口/函数都可以. 但是因为 存在了 类型约束, 所以此时该接口只能用来作为泛型约束参数
}

type MyNumber int

func (m MyNumber) String() string {
	return strconv.Itoa(int(m))
}

type MyString string

func (m MyString) String() string {
	return string(m)
}

func RangeConstraints[T MyGeneric](src []T) {
	Range(src)  // 因为 泛型参数是 []MyGeneric, 所以其本身也实现了 fmt.Stringer 接口, 约束放宽是可以的
	Range2(src) // 任意类型皆可
}
