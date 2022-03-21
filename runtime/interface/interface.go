package inter

type Duck interface {
	Quack()
}

type Cat struct {
	_ int
}

func (c *Cat) Quack() {}

func Quack() {
	var c Cat
	var d Duck = &c
	d.Quack()
}
