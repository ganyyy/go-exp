package meta

type IValue[T any] interface {
	NewValue() IValue[T]
	FromProto(T)   // From sets the value from the target.
	ToProto() T    // To gets the target from the value.
	Dyeing(func()) // Dyeing sets the dyeing root function.
}

type transfer[V, T any] struct {
	t2v   func(T) V
	v2t   func(V) T
	onSet func(V)
	onDel func(V)
}

// setHook is a hook that will be called when any field is set.
func (t *transfer[V, T]) setHook(v V) {
	if t.onSet != nil {
		t.onSet(v)
	}
}

// delHook is a hook that will be called when any field is del.
func (t *transfer[V, T]) delHook(v V) {
	if t.onDel != nil {
		t.onDel(v)
	}
}

type IDyeing interface {
	GetDyeing() func()
}

func v2v[V any](v V) V { return v }

func ValueTransfer[V any]() transfer[V, V] {
	return transfer[V, V]{
		t2v: v2v[V],
		v2t: v2v[V],
	}
}

func ReferenceTransfer[V IValue[T], T any](get IDyeing) transfer[V, T] {
	var transfer transfer[V, T]
	transfer.t2v = func(t T) (v V) {
		v = v.NewValue().(V)
		v.FromProto(t)
		return
	}
	transfer.v2t = func(v V) T { return v.ToProto() }
	transfer.onSet = func(v V) { v.Dyeing(get.GetDyeing()) }
	transfer.onDel = func(v V) { v.Dyeing(nil) }
	return transfer
}

func Pointer[T any](v T) *T { return &v }
