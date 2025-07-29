package main

import (
	"flag"
	"fmt"
	"os"

	"ganyyy.com/go-exp/demo/hotfix/enctext/sym"
)

var Key = []byte("123456")

func main() {
	var input, prefix string
	flag.StringVar(&input, "input", "", "input dynamic library path")
	flag.StringVar(&prefix, "prefix", "", "prefix for function names")
	flag.Parse()

	if input == "" || prefix == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	data, encFunc, err := sym.EncText(data, sym.GenRC4Enc(data, Key), prefix)
	if err != nil {
		panic(err)
	}

	if len(encFunc) == 0 {
		fmt.Println("No functions found to encrypt with prefix:", prefix)
		return
	}

	if err := os.WriteFile(input, data, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(input+".func", encFunc, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Encryption complete, modified file saved.")
}
