package main

import "unsafe"

func main() {

	var id = 100
	B(id)
}

type Map struct {
	M *struct {
		d map[int]*Data
	}
}

var m Map

func (m *Map) get(v int) *Data {
	return m.M.d[v]
}

var dataMap = make(map[int]*Data)

type Data struct {
	AA, AB, AC, AD, AE *int
}

//go:noinline
func B(id int) {
	var ptr int
	defer func() {

		println(recover())
		var buffer [1 << 20]int
		var ret int
		for i := range buffer {
			buffer[i] += i + 10
			ret += buffer[i]
		}
		println(ret)
	}()

	var pptr = unsafe.Pointer(&ptr)

	var data *Data
	println(uintptr(pptr))

	data = m.get(id)
	if data != nil {
		return
	}

	var p1, p2, p3, p4, p5 *int
	var a, b, c int
	var e int

	var f int

	p5 = &f
	println(p5)

	p1 = &a
	data = &Data{}
	p5 = &e
	*p5++

	p2 = &b
	p3 = &c

	data.AA = p1
	data.AB = p2
	data.AC = p3
	*p4 += *p5
	data.AD = p4
	println(pptr)
	m.M.d[id] = data
}

//go:noinline
func A1() *int {
	var a int
	var pa = &a
	return pa
}

//go:noinline
func A2() (*int, *int) {
	var a int
	var b int
	var pa = &a
	var pb = &b
	return pa, pb
}
