package main

import (
	"testing"
)

func TestFilter(t *testing.T) {
	genPath = append(genPath, "./demo")
	*typeNames = "Student,SomeData"
	*output = "./demo/filter.go"
	filter()
}
