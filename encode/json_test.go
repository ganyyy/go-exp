package main

import (
	"encoding/json"
	"testing"
)

type Type1 struct {
	Name string
	Addr string
	Age  int
}

type Type2 struct {
	Name string
	Age  int
}

func Test(t *testing.T) {
	a := Type2{
		Name: "123",
		Age:  100,
	}

	bs, _ := json.Marshal(&a)

	var b Type1

	_ = json.Unmarshal(bs, &b) // 不加omitempty, 也不影响 小转大
	t.Logf("%+v", b)
}
