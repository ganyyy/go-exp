package main

type iData struct {
	A int
	B int
	C int
	D int
}

func (d *iData) SetA(a int) { d.A = a }

//go:noinline
func (d *iData) SetB(b int)  { d.B = b }
func (d *iData) SetC(c int)  { d.C = c }
func (d *iData) SetD(d1 int) { d.D = d1 }

func (d *iData) GetA() int { return d.A }
func (d *iData) GetB() int { return d.B }
func (d *iData) GetC() int { return d.C }
func (d *iData) GetD() int { return d.D }
