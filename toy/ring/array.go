package ring

type Ring[T any] interface {
	Get(int) (T, bool)  // 获取指定位置的值, 并返回是否存在. 0表示正向第一个, -1表示逆向第一个
	Add(T)              // 添加一个新元素
	Range(func(T) bool) // 按照添加顺序迭代缓冲区
	Clear()             // 重置缓冲区
	Copy() []T          // 按照添加顺序copy缓冲区(浅拷贝)
	Len() int           // 返回实际长度
	Top() (T, bool)     // 获取首个元素, 等同于 Get(0)
	Pop() (T, bool)     // 弹出首个元素, 返回元素和弹出是否成功
}

const (
	initIndex = -1
)

func NewArrayRing[T any](capacity int32) *ArrayRing[T] {
	if capacity <= 0 {
		panic("invalid capicity")
	}
	var arr = &ArrayRing[T]{
		buffer: make([]T, capacity),
	}
	arr.resetIndex()
	return arr
}

type ArrayRing[T any] struct {
	read   int // 读指针, 表示读取的位置, 0开始.
	write  int // 写指针, 最后一次写入的位置, 起始值是initIndex
	buffer []T
}

// Add 添加一个新元素 批量插入的性能较差(耗时多一倍), 还是用
func (a *ArrayRing[T]) Add(ele T) {
	// if len(ele) <= 0 {
	// 	return
	// }
	var write = a.write
	var read = a.read
	var ln = len(a.buffer)
	// if len(ele) >= ln {
	// 	ele = ele[len(ele)-ln:]
	// }
	// for _, e := range ele {
	// 首个位置需要特殊处理一下
	if write != initIndex && read == write {
		read = (read + 1) % ln
	}
	write = (write + 1) % ln
	a.buffer[write] = ele
	// }
	a.read = read
	a.write = write
}

// Clear 重置整个ring
func (a *ArrayRing[T]) Clear() {
	if a.write == initIndex {
		return
	}
	// for gc
	// 这样的性能会更好点?
	var ele T
	for i := range a.buffer {
		a.buffer[i] = ele
	}
	a.resetIndex()
}

func (a *ArrayRing[T]) resetIndex() {
	a.read = 0
	a.write = initIndex
}

// Copy 浅拷贝元素的切片
func (a *ArrayRing[T]) Copy() []T {
	var ret = make([]T, 0, a.Len())
	a.Range(func(t T) bool {
		ret = append(ret, t)
		return true
	})
	return ret
}

// Get 获取指定位置的元素, -1表示最后一个, 0表示第一个
func (a *ArrayRing[T]) Get(idx int) (T, bool) {
	var ln = a.Len()
	if idx < 0 {
		idx = ln + idx
	}
	if idx < 0 || idx >= a.Len() {
		var ele T
		return ele, false
	}
	return a.buffer[(a.read+idx)%ln], true
}

// Range 迭代切片, 返回false则不会继续迭代
func (a *ArrayRing[T]) Range(cb func(T) bool) {
	if cb == nil || a.Len() <= 0 {
		return
	}
	// 迭代, 返回是否继续
	var iter = func(buf []T) bool {
		for _, v := range buf {
			if !cb(v) {
				// 函数中止, 停止迭代
				return false
			}
		}
		// 继续迭代
		return true
	}
	read := a.read
	write := a.write
	if read <= write {
		// 还未走完一圈
		iter(a.buffer[read : write+1])
		return
	}
	// [read, ln]
	if !iter(a.buffer[read:]) {
		return
	}
	// [0, write]
	iter(a.buffer[:write+1])
}

func (a *ArrayRing[T]) Len() int {
	write := a.write
	if write == initIndex {
		return 0
	}
	read := a.read
	ln := a.bufferLen()
	if read <= write {
		// [read, write]
		return write + 1 - read
	}
	// [read, ln] + [0, write]
	return ln - read + write + 1
}

func (a *ArrayRing[T]) bufferLen() int {
	return len(a.buffer)
}

func (a *ArrayRing[T]) Top() (ele T, ok bool) {
	if a.Len() == 0 {
		a.Clear()
		return
	}
	return a.Get(0)
}

func (a *ArrayRing[T]) Pop() (ele T, ok bool) {
	ele, ok = a.Top()
	if !ok {
		return
	}
	var empty T
	read := a.read
	a.buffer[read] = empty // for GC
	if read == a.write {
		// 空了, 重置一下
		a.resetIndex()
		return
	}
	a.read = (read + 1) % a.bufferLen()
	return
}
