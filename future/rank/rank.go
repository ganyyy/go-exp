package rank

type Compare interface {
	Compare() int // -1 小于, 0 等于, 1 大于
}

type LimitRank[K comparable, V Compare] struct {
	elements    map[K]V
	size        int
	ElementList []V
}

func (r *LimitRank[K, V]) Init(size int) {
	r.elements = make(map[K]V, size)
	r.size = size
}

func (r *LimitRank[K, V]) GetElement(k K) (V, bool) {
	var v, ok = r.elements[k]
	return v, ok
}

func (r *LimitRank[K, V]) GetElementWithDefault(k K, def V) V {
	var v, ok = r.GetElement(k)
	if !ok {
		return def
	}
	return v
}

func (r *LimitRank[K, V]) GetRank() int {
	return 0
}
