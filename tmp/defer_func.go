package main

import (
	"fmt"
)

func DeferFunc4() (t int) {
	t = 1
	defer func(i int) {
		fmt.Println(i, &i)
		fmt.Println(t, &t)
	}(t)
	return 2
}

func main() {
	DeferFunc4()
}
