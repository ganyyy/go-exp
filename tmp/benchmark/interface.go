package main

type Task interface {
	Do()
}

type TaskOne interface {
	Task
	GetA() int
}

type TaskTow interface {
	Task
	GetB() int
	SetB(int)
}

type Task3 interface {
	Task
	SetA(int)
	SetB(int)
	SetC(int)
	SetD(int)
}

type Task4 interface {
	Task
	GetA()
	GetB()
	GetC()
	GetD()
	Do2(int)
}

type TaskBase struct{}

func (TaskBase) Do() {}

type AA struct{}

func (AA) GetA() int { return 0 }
func (AA) SetA(int)  {}

type BB struct{}

func (BB) GetB() int { return 0 }
func (BB) SetB(int)  {}

type CC struct{}

func (CC) GetC() int { return 0 }
func (CC) SetC(int)  {}

type DD struct{}

func (DD) GetD() int { return 0 }
func (DD) SetD(int)  {}

type TTOne struct {
	TaskBase
	AA
}

type TTOne2 struct{ TTOne }
type TTOne3 struct{ TTOne }
type TTOne4 struct{ TTOne }
type TTOne5 struct{ TTOne }

type TTTow struct {
	TaskBase
	BB
}

type TTTow2 struct{ TTTow }
type TTTow3 struct{ TTTow }
type TTTow4 struct{ TTTow }
type TTTow5 struct{ TTTow }

type TT3 struct {
	TaskBase
	AA
	BB
	CC
	DD
}

type TT32 struct{ TT3 }
type TT33 struct{ TT3 }
type TT34 struct{ TT3 }
type TT35 struct{ TT3 }

type TT4 struct {
	TaskBase
	AA
	BB
	CC
	DD
}

type TT42 struct{ TT4 }
type TT43 struct{ TT4 }
type TT44 struct{ TT4 }
type TT45 struct{ TT4 }

func (TT4) Do2(int) {}

func DoTask(t Task) {
	switch run := t.(type) {
	case TaskOne:
		run.GetA()
		run.Do()
	case TaskTow:
		run.SetB(100)
		run.Do()
	case Task3:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case Task4:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	}
}

func DoTask2(t Task) {
	switch run := t.(type) {
	case TTOne:
		run.GetA()
		run.Do()
	case TTOne2:
		run.GetA()
		run.Do()
	case TTOne3:
		run.GetA()
		run.Do()
	case TTOne4:
		run.GetA()
		run.Do()
	case TTOne5:
		run.GetA()
		run.Do()
	case TTTow:
		run.SetB(100)
		run.Do()
	case TTTow2:
		run.SetB(100)
		run.Do()
	case TTTow3:
		run.SetB(100)
		run.Do()
	case TTTow4:
		run.SetB(100)
		run.Do()
	case TTTow5:
		run.SetB(100)
		run.Do()
	case TT3:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case TT32:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case TT33:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case TT34:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case TT35:
		run.SetA(0)
		run.SetB(0)
		run.SetC(0)
		run.SetD(0)
		run.Do()
	case TT4:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	case TT42:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	case TT43:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	case TT44:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	case TT45:
		run.GetA()
		run.GetB()
		run.GetC()
		run.GetD()
		run.Do()
		run.Do2(100)
	}
}

var arr = []Task{
	TTOne{},
	TTTow{},
	TTOne{},
	TTTow{},
	TTOne{},
	TTTow{},
	TTOne{},
	TTTow{},
	TTOne{},
	TT3{},
	TT4{},
	TTOne2{},
	TTOne2{},
	TTOne3{},
	TTOne3{},
	TTOne4{},
	TTTow2{},
	TT32{},
	TT42{},
	TTOne4{},
	TTTow3{},
	TT33{},
	TT43{},
	TTOne5{},
	TTTow4{},
	TT34{},
	TT44{},
	TTOne5{},
	TTTow5{},
	TT35{},
	TT45{},
}

func main() {

	for _, t := range arr {
		DoTask(t)
		DoTask2(t)
	}
}
