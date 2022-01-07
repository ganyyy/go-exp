package main

type TestStruct struct {
	Val  *byte
	Next int
}

func largeAlloc() {

	showTypeInfo(TestStruct{})

	var tmp = make([]TestStruct, 100000)

	_ = tmp
}
