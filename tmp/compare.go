package main

type CheckType struct {
	V1, V2, V3 int
}

//go:noinline
func CompareCheck1(c1, c2 CheckType) bool {

	if c1.V1 != c2.V1 {
		return c1.V1 > c2.V1
	}

	if c1.V2 != c2.V2 {
		return c1.V2 > c2.V2
	}

	if c1.V3 != c2.V3 {
		return c1.V3 > c2.V3
	}

	return false
}

//go:noinline
func CompareCheck2(c1, c2 CheckType) bool {

	if c1.V1 > c2.V1 {
		return true
	} else {
		if c1.V2 > c2.V2 {
			return true
		} else {
			if c1.V3 > c2.V3 {
				return true
			}
		}
	}

	return false
}
