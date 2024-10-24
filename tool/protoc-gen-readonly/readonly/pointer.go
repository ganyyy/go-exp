package readonly

type Pointer[T any] struct {
	inner *T
}

func NewPointer[T any](inner **T) Pointer[T] {
	if inner == nil || *inner == nil {
		return Pointer[T]{inner: nil}
	}
	p := *inner
	*inner = nil
	return Pointer[T]{inner: p}
}

func (x *Pointer[T]) IsNil() bool {
	return x == nil || x.inner == nil
}

func (x *Pointer[T]) Get() T {
	if x.IsNil() {
		var empty T
		return empty
	}
	return *x.inner
}
