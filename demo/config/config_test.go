package main

import (
	"encoding/json"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	type UnSupportJson struct {
		Chan chan int
		F    func(int)
	}

	var p UnSupportJson

	_, err := json.Marshal(p)
	t.Logf("err:%v", err)

}
