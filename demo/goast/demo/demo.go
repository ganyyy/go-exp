package demo

type SomeData struct {
	Address string
}

//go:generate goast -type Student
type Student struct {
	SomeData *SomeData `json:"data" bson:"some_data2323"`
	Name     string    `json:"name" bson:"Name"`
	Age      int       `json:"age" bson:"Age2"`
	Other    string    `bson:"ooo"`
	Other2   string    `bson:"ooo3"`
	Other3   string    `bson:"aaaaaa"`
}

func AAA(a int) int {
	t := [100]int{}
	return t[99] + a
}

func BBB() {
	AAA(10)
}
