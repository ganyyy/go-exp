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
