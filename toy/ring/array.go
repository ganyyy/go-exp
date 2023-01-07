package ring

type Ring[T any] interface {
	Get(int) (T, bool)
	Add(T)
	Range(func(T) bool)
	Clear()
	Copy() []T
	Len() int
}

const (
	initIndex = -1
)

func NewArrayRing[T any](capacity int32) Ring[T] {
	if capacity <= 0 {
		panic("invalid capicity")
	}
	return &arrayRing[T]{
		buffer: make([]T, capacity),
		read:   initIndex,
		write:  initIndex,
	}
}

type arrayRing[T any] struct {
	read   int // 读指针, 表示读取的位置, 起始值是initIndex, write写入一圈后开始增加
	write  int // 写指针, 最后一次写入的位置, 起始值是initIndex
	buffer []T
}

// Add 添加一个新元素
func (a *arrayRing[T]) Add(ele T) {
	var write = a.write
	var read = a.read
	var ln = len(a.buffer)
	var over = write+1 >= ln
	write = (write + 1) % ln
	a.buffer[write] = ele
	if over && read == initIndex {
		read = 0
	}
	if read == write {
		a.read = (read + 1) % ln
	}
	a.write = write
}

// Clear 重置整个ring
func (a *arrayRing[T]) Clear() {
	var ele T
	for i := range a.buffer {
		a.buffer[i] = ele
	}
	a.read = initIndex
	a.write = initIndex
}

// Copy 浅拷贝元素的切片
func (a *arrayRing[T]) Copy() []T {
	var ret = make([]T, 0, a.Len())
	a.Range(func(t T) bool {
		ret = append(ret, t)
		return true
	})
	return ret
}

// Get 获取指定位置的元素, -1表示最后一个, 0表示第一个
func (a *arrayRing[T]) Get(idx int) (T, bool) {
	var ln = a.Len()
	if idx < 0 {
		idx = ln + idx
	}
	if idx < 0 || idx >= a.Len() {
		var ele T
		return ele, false
	}
	var read = a.read
	if read == initIndex {
		read = 0
	}
	return a.buffer[(read+idx)%ln], true
}

// Range 迭代切片, 返回false则不会继续迭代
func (a *arrayRing[T]) Range(cb func(T) bool) {
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
	if read == initIndex {
		// [read, write]
		iter(a.buffer[0 : a.write+1])
	} else {
		// 这一定是走了一圈了
		// [read, ln], [0, read]
		if !iter(a.buffer[read:]) {
			return
		}
		iter(a.buffer[:read])
	}
}

func (a *arrayRing[T]) Len() int {
	write := a.write
	if write == initIndex {
		return 0
	}
	read := a.read
	if read == initIndex {
		return a.write + 1
	}
	return len(a.buffer)
}
