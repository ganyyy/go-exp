//go:build ignore

package main

import (
	"fmt"
)

type TestNoCopy struct {
}

func switchFunc(a int) int {
	switch a {
	case 0:
		return 100
	case 1:
		return 200
	case 2:
		return 700
	case 3:
		return 1900
	case 4:
		return 10230
	case 5:
		return 100123
	default:
		return 300
	}
}

func ifFunc(a int) int {
	if a == 0 {
		return 100
	} else if a == 1 {
		return 200
	} else if a == 2 {
		return 700
	} else if a == 3 {
		return 1900
	} else if a == 4 {
		return 10230
	} else if a == 5 {
		return 100123
	} else {
		return 300
	}
}

func main() {
	var v int
	_, _ = fmt.Scanf("%d", &v)
	println(v)
	println(switchFunc(v), ifFunc(v))
}
