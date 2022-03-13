package helper

func To[T any](v any) T {
	var ret, ok = v.(T)
	if !ok {
		var rr T
		return rr
	}
	return ret
}
