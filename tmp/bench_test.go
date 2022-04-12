package main

import (
	"strconv"
	"testing"
)

func BenchmarkCompareCheckType(b *testing.B) {

	var src = [][2]CheckType{
		{
			{
				V1: 1,
				V2: 0,
				V3: 0,
			},
			{
				V1: 0,
				V2: 0,
				V3: 0,
			},
		},
		{
			{
				V1: 0,
				V2: 1,
				V3: 0,
			},
			{
				V1: 0,
				V2: 0,
				V3: 0,
			},
		},
		{
			{
				V1: 0,
				V2: 0,
				V3: 1,
			},
			{
				V1: 0,
				V2: 0,
				V3: 0,
			},
		},
	}

	for i, t := range src {
		b.Run("Check1"+strconv.Itoa(i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CompareCheck1(t[0], t[1])
			}
		})

		b.Run("Check2"+strconv.Itoa(i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CompareCheck2(t[0], t[1])
			}
		})
	}

}

func TestChannel1(t *testing.T) {
	var ch chan int

	ch <- 1
}

func TestChannel2(t *testing.T) {
	var ch = make(chan int)

	close(ch)

	ch <- 1
}
