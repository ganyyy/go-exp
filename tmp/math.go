//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math"
)

func main() {
	v := math.Log(-1)

	var m = map[float64]int{
		v: 10,
		v: 20,
		v: 30,
	}

	fmt.Println(m[v], len(m))
}
