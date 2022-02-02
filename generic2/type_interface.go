package generic2

func Map[F, T any](src []F, f func(F) T) []T {
	var ret = make([]T, 0, len(src))
	for _, s := range src {
		ret = append(ret, f(s))
	}
	return ret
}

func Filter[F any](src []F, f func(F) bool) []F {
	// 最多扩容一次
	var ret = make([]F, 0, len(src)>>1)
	for _, s := range src {
		if !f(s) {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func Reduce[F any](src []F, f func(a, b F) F, d ...F) F {
	var def F
	if len(d) >= 1 {
		def = d[0]
	}
	if len(src) == 0 {
		return def
	}
	var cur = f(def, src[0])
	for _, s := range src[1:] {
		cur = f(cur, s)
	}
	return cur
}

// 关于参数的匹配规则

func ShowTypeInterface1(_ []map[int]bool)                     {}
func ShowTypeInterface2[T1 any](_ T1)                         {}
func ShowTypeInterface3[T1 any](_ []T1)                       {}
func ShowTypeInterface4[T1 comparable, T2 any](_ []map[T1]T2) {}

/*
	1. 自动推导仅限于泛型函数, 泛型类型还是要声明具体的约束
	2. 泛型函数的自动推断的前提是参数有泛型参数. 如果仅用于函数体或者返回值会无法自动推断
*/

/*
类型统一
	1. 忽略所有的 untyped const 类型的实参和其对应的形参泛型参数
	2. 匹配有类型的实参和泛型参数, 并反过来确定无类型的常量对应的泛型类型是否符合有类型实参的预期
       如果都是无类型的常量, 会使用其默认类型. 比如 1 的默认类型是 int, 1.0的默认类型是float64
*/
