package main

import (
	"log"
	"testing"
	"time"
)

func RangeSlice(src *[]int) {
	for i, v := range *src {
		log.Printf("i:%v, v:%v", i, v)
		time.Sleep(time.Second)
	}
}

func IndexSlice(src *[]int) {
	defer func() {
		log.Printf("error:%v", recover())
	}()
	for i := range *src {
		log.Printf("i:%v, v:%v", i, (*src)[i])
		time.Sleep(time.Second)
	}
}

func TestReslice(t *testing.T) {
	var src = []int{1, 2, 3, 4, 5}

	go func() {
		time.Sleep(time.Second)
		src = src[:3]
	}()
	go RangeSlice(&src)
	go IndexSlice(&src)

	time.Sleep(time.Second * time.Duration(len(src)))
}

func TestNilStruct(t *testing.T) {

	func() {
		defer func() {
			t.Log(recover())
		}()
		var n *NilStruct2
		t.Logf("%+v", n.GetName())
	}()

	func() {
		var n = new(NilStruct2)
		t.Logf(n.GetName())
	}()

}
