package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	var f, err = os.Open("/Users/Gann/Downloads/39562.txt")
	if err != nil {
		panic(err)
	}
	var bs, _ = io.ReadAll(f)

	fmt.Println(string(bs))

	var src map[interface{}]interface{}

	err = json.Unmarshal(bs, &src)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range src {
		fmt.Println(k, v)
	}
}
