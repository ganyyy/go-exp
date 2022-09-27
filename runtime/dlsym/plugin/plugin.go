package main

import "ganyyy.com/go-exp/runtime/dlsym/pkg"

const (
	AddName  = "Add"
	Add2Name = "Add2"
)

//go:noinline
func Add(a, b int) int {
	return a + b
}

//go:noinline
func Add2(a, b int) int {
	return a + b + 100
}

func FuncAddr() map[string]uintptr {
	return map[string]uintptr{
		pkg.FuncName(pkg.Add): pkg.FuncEntry(pkg.Add),
		AddName:               pkg.FuncEntry(Add),
		Add2Name:              pkg.FuncEntry(Add2),
	}
}
