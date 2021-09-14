package main

import (
	"encoding/json"
	"fmt"
)

type TestPtr struct {
	Name string
}

type TestVal struct {
	Val *TestPtr `json:"val,omitempty"`
}

func main() {
	var val = TestVal{}

	var bs, _ = json.Marshal(val)
	fmt.Println(string(bs))

	val = TestVal{
		Val: &TestPtr{
			Name: "213",
		},
	}

	_ = json.Unmarshal(bs, &val)
	fmt.Println(val)
}
