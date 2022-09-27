package main

import (
	"testing"
)

func TestFuncName(t *testing.T) {

	t.Log(logName(Data))
	t.Log(logName(Stu.Name1))
	t.Log(logName((*Stu).Name2))
}
