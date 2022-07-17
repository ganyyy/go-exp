package demo

import "ganyyy.com/go-exp/demo/goast/demo2"

type SomeData struct {
	Address string `bson:"addr"`
}

type Generic[T any] struct {
	V T
}

type Generic2[T, T2, T3 any] struct {
	V  T
	V2 T2
	V3 T3
}

//go:generate goast -type Student
type Student struct {
	SomeData  *SomeData              `json:"data" bson:"some_data2323"`
	SomeData2 SomeData               `json:"data3" bson:"some_data23233"`
	SomeData4 ****SomeData           `json:"data4" bson:"some_data23234"`
	Name      string                 `json:"name" bson:"Name"`
	Age       int                    `json:"age" bson:"Age2"`
	Other     string                 `bson:"ooo"`
	Other2    string                 `bson:"ooo3"`
	Other3    string                 `bson:"aaaaaa"`
	Invalid   chan int               `bson:"invalid"`
	Map       map[*SomeData]SomeData `bson:"mm"`
	TTT       [][][][]int            `bson:"tt"`
	Omit      int                    `bson:",omitempty"`
	Generic   Generic[int]           `bson:"ge"`
	Generic2  Generic2[int, int, string]
	Generic3  Generic2[int, int, Generic[int]]
	T1        demo2.Lala
	T2        demo2.GLala[int]
}

func AAA(a int) int {
	t := [100]int{}
	return t[99] + a
}

func BBB() {
	AAA(10)
}
