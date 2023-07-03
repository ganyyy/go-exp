package main

import (
	"flag"
	"fmt"
	"os"
)

func addSumData(s string, a, b int) int {
	fmt.Println(s)
	for i := 0; i < a; i++ {
		b += i * a
	}
	return b
}

var s = flag.String("s", "hello", "string flag")

func main() {
	flag.Parse()
	fmt.Println("hello world")
	fmt.Println(addSumData(*s, 10, 20))

	f, e := os.OpenFile("./test.log", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	for i := 0; i < 10; i++ {
		f.WriteString(fmt.Sprintf("%s %d\n", *s, i))
	}
}
