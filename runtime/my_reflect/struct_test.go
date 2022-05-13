package myreflect

import (
	"encoding/json"
	"testing"
)

type V interface{}

func TestCreateStruct(t *testing.T) {
	type Value struct {
		V V
		A int
		B int
	}

	var m = map[string]int{
		"V1":      10,
		"V2":      20,
		"V3":      30,
		"V4":      40,
		"V5":      50,
		"Sadasd":  11,
		"Sadasd2": 11,
		"Sadasd3": 11,
		"Sadasd4": 11,
	}

	var vv Value
	vv.V = CreateStruct(m).Interface()

	var bs, _ = json.Marshal(vv)

	t.Logf("%+v, %v", vv.V, string(bs))

	bs, _ = json.Marshal(m)
	t.Logf("%+v", string(bs))
}
