package generic2

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func Double[E constraints.Integer](s []E) []E {
	var ret = make([]E, len(s))
	for i, v := range s {
		ret[i] = v + v
	}
	return ret
}

type Slice[E any] interface {
	~[]E
}

func DoubleDefined[E constraints.Integer, S Slice[E]](s S) S {
	var ret = make([]E, len([]E(s)))
	for i, v := range s {
		ret[i] = v + v
	}
	return S(ret)
}

func DoubleDefined2[E constraints.Integer, S ~[]E](s S) S {
	var ret = make([]E, len([]E(s)))
	for i, v := range s {
		ret[i] = v + v
	}
	return S(ret)
}

/*
带有指针的方法
*/

/*
当前这种声明模式存在以下的问题
1. Set的接收器是指针类型
2. 指针类型的切片其初始化都是空指针. 运行时panic
*/

type Setter interface {
	Set(string)
}

type Settable int

func (s *Settable) Set(v string) {
	var n, _ = strconv.Atoi(v)
	*s = Settable(n)
}

//FromStrings 这种情况就是泛型参数只存在于函数体/返回值中的情况
//所以无法自动推断
func FromStrings[T Setter](s []string) []T {
	var ret = make([]T, len(s))
	for i, v := range s {
		ret[i].Set(v)
	}
	return ret

}

//Setter2 约束了 *T必须要有Set(string)方法
type Setter2[B any] interface {
	Set(string)
	*B
}

func FromStrings2[T any, PT Setter2[T]](s []string) []T {
	// 切片元素类型是T, 而不是*T, 避免了空指针的问题
	var ret = make([]T, len(s))
	for i, v := range s {
		p := PT(&ret[i]) // 卧槽, 还能这么用?
		p.Set(v)
	}
	return ret
}

// 没有运算符重载又如何? 没啥是不能包装一层解决的

//NumberEqual 这个要怎么理解呢?
//拥有一个和T对比返回bool的类型约束
type NumberEqual[T any] interface {
	Equal(T) bool
}

type EqualNum int

func (e EqualNum) Equal(v EqualNum) bool {
	return int(e) == int(v)
}

//IndexEqual 泛型约束引用自身. 这里的T恰好是本身
//说明需要 T.Equal(T) bool
func IndexEqual[T NumberEqual[T]](src []T, e T) int {
	for i, v := range src {
		if v.Equal(e) {
			return i
		}
	}
	return -1
}

//IndexEqual2 简单点~
func IndexEqual2[T interface{ Equal(T) bool }](src []T, e T) int {
	for i, v := range src {
		if v.Equal(e) {
			return i
		}
	}
	return -1
}

/*
泛型不是接口. 接口总是依赖于指针, 这就意味着会出现装箱操作
泛型本质上是编译期进行的特殊处理, 所以不会装箱
*/
