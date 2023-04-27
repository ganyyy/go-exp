package main

type Data struct {
	A int
	B int
	C int
	D int
}

func (d *Data) GetA() int { return d.A }
func (d *Data) GetB() int { return d.B }
func (d *Data) GetC() int { return d.C }
func (d *Data) GetD() int { return d.D }

func (d *Data) SetA(a int) { d.A = a }

//go:noinline
func (d *Data) SetB(b int)  { d.B = b }
func (d *Data) SetC(c int)  { d.C = c }
func (d *Data) SetD(d1 int) { d.D = d1 }
