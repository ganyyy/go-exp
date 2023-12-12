package generic2

/*
	约束中的复合类型, 会尽量提取相同的部分.
	(如果其中不包含共同的部分, 那么本质上这个约束是无法实例化的)错误的

	对于一个复合约束而言, 对其进行操作的前提是: 二者的某次的输入/输出的类型一致

	即使无法直接使用相关的操作, 也可以约束输入的参数, 并通过断言的形式处理逻辑
*/

// ByteSeq string和[]byte的共同点是:
// 允许索引, 拥有len, ...?
type ByteSeq interface {
	~string | ~[]byte
}

func Join[T ByteSeq](src []T, seq T) (ret T) {
	if len(src) == 0 {
		return
	}
	if len(src) == 1 {
		return T(append([]byte(nil), []byte(src[0])...))
	}

	var bSeq = []byte(seq)

	// 计算总的长度
	var n = (len(src) - 1) * len(bSeq)
	for _, v := range src {
		n += len([]byte(v))
	}
	var b = make([]byte, n)

	// 依次拷贝赋值
	bp := copy(b, []byte(src[0]))
	for _, s := range src[1:] {
		bp += copy(b[bp:], bSeq)
		bp += copy(b[bp:], []byte(s))
	}
	return T(b)
}

// StructField 在这个复合类型中, 都拥有的属性是 x
// 但是其类型不一致, 只能使用断言的方式来进行操作
type StructField interface {
	struct {
		a int
		x int
	} | struct {
		b int
		x float64
	} | struct {
		c int
		x bool
	}
}

// IncrementX 这种方式只能用断言了..
func IncrementX[T StructField](p *T) {
	switch v := interface{}(p).(type) {
	case *struct {
		a int
		x int
	}:
		v.x += 100
	case *struct {
		b int
		x float64
	}:
		v.x *= 10
	case *struct {
		c int
		x bool
	}:
		v.x = false
	default:
		println("unknown")
	}
}

type IntSliceOrMap interface {
	[]int | map[int]int // 这里不能使用底层类型, 否则无法断言...
}

// Entry 这个例子也不行. 原则上输出/输出都是一致的, 为啥编译不过呢?
func Entry[T IntSliceOrMap](c T, i int) int {
	switch t := interface{}(c).(type) {
	case []int: // 期望后期会推出  ~[]int 这种形式的断言方法
		return t[i]
	case map[int]int:
		return t[i]
	default:
		return 0
	}
}

/*
	多重类型约束
*/

type SliceType[T any] interface {
	~[]T
}

// MapSlice 这个的有点是可以保持原有的类型!
func MapSlice[S SliceType[E], E any](s S, f func(E) E) S {
	var r = make([]E, len([]E(s)))
	for i, v := range s {
		r[i] = f(v)
	}
	return S(r)
}
