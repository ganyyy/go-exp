package main

import (
	"fmt"

	"mydomain.com/foo"
)

func Bar() {
	foo.Foo()
	fmt.Println("In Bar package")
}
