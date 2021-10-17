//go:build ignore
// +build ignore

package main

func printArgs(a int) (a2 int) {
	return a
}

func main() {
	defer printArgs(1)
}
