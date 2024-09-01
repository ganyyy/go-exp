package meta

import "reflect"

func Equal[T comparable](a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if !ta.Comparable() {
		return false
	}
	return a == b
}
