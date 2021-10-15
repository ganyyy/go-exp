package package_b

import (
	"log"
)

type inner struct {
	name string
	age int
	tmp [1<<20]int
}

func (i *inner) show() {
	log.Printf("%+v", *i)
}

func GetInner(name string, age int) *inner {
	var in = &inner{
		name: name,
		age:  age,
	}
	in.tmp[0] = 100
	return in
}

type Show interface {
	show()
}

type Outer struct {
	Name string
	Age int
}