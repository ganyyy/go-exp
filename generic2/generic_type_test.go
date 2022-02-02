package generic2

import (
	"fmt"
	"testing"
)

func TestGenericType(t *testing.T) {
	// 声明时, 必须添加具体的类型, 不能自动推导(?)
	_ = MyVector[int]{1, 2, 3, 4}
	_ = MyVector[string]{"1", "2", "3,", "4"}
	_ = MyVector[MyNumber]{MyNumber(1), MyNumber(100)}

	// 支持使用接口来作为类型约束, 但是不能使用泛型约束(MyGeneric 就不行)
	_ = MyVector[fmt.Stringer]{MyNumber(1), MyString("123123")}

	// 对于泛型类型引用自身而言, 顺序不重要
	var p = P[int, string]{
		V1: 100,
		V2: "str",
	}
	var f P[int, string]
	var pp P[string, int]
	p.PP = &pp
	p.F = &f

	var ppp = P[int, int]{}
	_ = ppp
	//ppp.PP = &p // 类型参数不匹配, 所以无法赋值

	t.Logf("%+v", p)
}
