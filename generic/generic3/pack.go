package generic3

func Pack[T any](t ...T) []T {
	return t
}

type Set[K comparable] map[K]struct{}

func PackSet[K comparable](keys ...K) Set[K] {
	var ret = make(map[K]struct{}, len(keys))
	for _, k := range keys {
		ret[k] = struct{}{}
	}
	return ret
}

func UnpackSet[K comparable, S ~map[K]struct{}](set S) []K {
	var ret = make([]K, 0, len(set))
	for k := range set {
		ret = append(ret, k)
	}
	return ret
}

type SomeGeneric[T any] interface {
	GetValue() T
}

type (
	IntGeneric = SomeGeneric[int]
)

type IntBase struct{ V int }

func (i *IntBase) GetValue() int { return i.V }

func DoGeneric[T any](v SomeGeneric[T]) T {
	return v.GetValue()
}

func DoIntGeneric(v IntGeneric) int {
	return v.GetValue()
}
