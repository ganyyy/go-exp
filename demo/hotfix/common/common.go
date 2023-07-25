package common

import (
	_ "unsafe"
)

type Data struct {
	A, B, C, D int
}

// SetA
//
//go:noinline
func (m *Data) SetA(a int) {
	m.A = a
	m.setC(a + 1)
}

// SetB
func (m *Data) SetB(b int) {
	m.B = b
}

// setC
//
//go:noinline
func (m *Data) setC(c int) {
	m.C = c
}
