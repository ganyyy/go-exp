package foo

import (
	"fmt"
	"ganyyy.com/go-exp/generic2"
)

func Foo() {
	fmt.Println("In Foo package")

	generic2.Index([]int(nil), 10)
}
