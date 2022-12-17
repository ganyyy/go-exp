package main

import "testing"

type InnerStruct struct {
	Data int
}

func (i *InnerStruct) GetData(t *testing.T) int {
	t.Logf("innert ptr:%p", i)
	return i.Data
}

type OutStruct struct {
	*InnerStruct
}

func ShowInfo(o *OutStruct, t *testing.T) {
	t.Logf("out ptr:%p", o)
	o.GetData(t)
}

func TestTestDataAddress(t *testing.T) {
	var o OutStruct
	o.InnerStruct = &InnerStruct{}
	ShowInfo(&o, t)
}
