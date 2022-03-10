package generic2

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

/*
	类型转换: 如果集合中的类型可以显式的转换, 那么就意味着这个泛型类型也可以转换
*/

func Convert[From, To interface {
	constraints.Integer | constraints.Float
}](f From) To {
	return To(f)
}

/*
	对类型集而言, 需要保证所有的非类型常量满足最低限度的要求.
	比如带有int8的集合中, 最大的整型常量不能超过127, 最小的整型常量不能低于-128
*/

func AddNum[T constraints.Integer](a T) T {
	return a + 1 // [-128, 127]
}

/*
	当一个约束中嵌入了多个约束(每个约束都是一行)时,
	最终的约束获取的是所有子约束的交集
*/

func AddString[T interface {
	constraints.Ordered // 交集只有~string
	~string
}](v T) T {
	return v + v
}

/*
	使用 | 链接的子约束,
	最终结果是取子约束的并集
*/

func AddUnion[T interface {
	constraints.Signed | constraints.Unsigned
}]() {

}

/*
	当类型和接口联合时, 将没有任何可以直接使用的操作. 可以使用接口断言
	(直接禁用了卧槽)
	所以接口一般是嵌入, 类型一般是联合
	但是可以用来进行类型检查(错误的)
*/

type StringFish interface {
	string
}

func RunStringFish[T StringFish](T) {

}

/*
	当约束中, 所有嵌入的约束无法同时满足时, 这个约束就无法实例化
*/

type Unsatisfiable interface {
	int | float32
	fmt.Stringer
}

/*
	所有的泛型类型在运行时都不存在, 反射只能获取其原始类型, 而无法获取泛型类型!
*/

/*
	其他好玩的问题
	1. switch type 支持 底层类型匹配
	2. 泛型方法对于duck type的语言不好实现
	3. 参是无法隐式的做值和地址的转换, 这就意味着方法的接收器决定了实参的类型
		- 比如: bytes.Buffer并没有实现String() string方法
*/
