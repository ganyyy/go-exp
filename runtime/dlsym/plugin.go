package main

//go:noinline
func Sum(i int) int {
	var ret int

	println(i)

	for j := 0; j < i; j++ {
		ret += j
	}
	return ret
}

func Show(i int) {
	Sum(i)
}

var ReplaceFunction = map[string]string{
	"Sum": "main.TotalSum",
}
