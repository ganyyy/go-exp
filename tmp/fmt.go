//go:build ignore

package main

import (
	"fmt"
)

type MyString struct {
	A, B, C int
}

func (m *MyString) String() string {
	if m == nil || m.A != 0 {
		return "nil"
	}
	return fmt.Sprintf("%v, %v, %v", m.A, m.B, m.C)
}

func main() {
	var m *MyString

	fmt.Printf("%s\n", m)
}
