package main

func Sum(i int) int {
	var ret int

	for j := 0; j < i; j++ {
		ret += j
	}
	return ret
}

var ReplaceFunction = map[string]string{
	"Sum": "main.TotalSum",
}
