package monkey_patch

//go:noinline
func Add(a, b int) int {
	return a + b
}

type Runnable struct {
	AAA int
}

//go:noinline
func (r *Runnable) SetAAA(a int) {
	r.AAA = a
}
