package enum

//go:generate goenum -type=Data,Data2
type Data int

const (
	D1 Data = iota
	D2
	D3
	D4
	_
	D6
	D7
	D8
)

type Data2 int

const (
	DD1 Data2 = iota
	DD2
	DD3
	DD4
)
