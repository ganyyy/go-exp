package generic2

//MyVector 声明泛型参数类型
type MyVector[T any] []T

//Add 泛型类型方法定义时, 必须要添加泛型参数声明
func (v *MyVector[T]) Add(ele T) {
	//当前IDE有点Bug, 无法识别泛型类型的原始类型
	var tmp = []T(*v)
	tmp = append(tmp, ele)
	*v = tmp
}

type MyMap[K comparable, V any] map[K]V

//P 泛型约束作为属性类型. 如果引用自身, 那么自身也需要使用泛型约束
type P[T1, T2 any] struct {
	F  *P[T1, T2] // 这两种类型声明都是合法的(!)
	PP *P[T2, T1] // 草案中说的是必须要保证类型参数的顺序一致性
	V1 T1
	V2 T2
}

//Add method must have no type parameters
//泛型类型不允许添加泛型方法
//func (t P[T1, T2]) Add[T3 any]() {
//
//}

//ListHead 间接引用的类型约束也会传递
//泛型的这种写法, 容易和切片搞混欸...
type ListHead[T any] struct {
	Head *ListElement[T]
}

type ListElement[T any] struct {
	Next *ListElement[T]
	Val  T
	Head *ListHead[T]
}

type Normal struct {
}

//Show 即使是非泛型类型, 也不允许添加类型参数方法
//func (n Normal) Show[T any](a T) {}
