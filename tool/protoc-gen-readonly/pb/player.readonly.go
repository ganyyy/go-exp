// Code generated by protoc-gen-readonly. DO NOT EDIT.
package pb

import (
	data "protoc-gen-readonly/pb/data"
	readonly "protoc-gen-readonly/readonly"
)

const _ = readonly.Version

type PlayerReadOnly struct {
	inner   *Player
	_Sa     readonly.List[*data.SimpleDataReadOnly]
	_Ma     readonly.Map[int32, *data.ReferencedDataReadOnly]
	_Da     *data.SimpleDataReadOnly
	_Age    readonly.Pointer[int32]
	_State2 readonly.Pointer[State]
	_Data   readonly.List[byte]
}

func NewPlayerReadOnly(p *Player) *PlayerReadOnly {
	if p == nil {
		return &PlayerReadOnly{inner: nil}
	}
	inner := p
	return &PlayerReadOnly{
		inner:   inner,
		_Sa:     readonly.NewListFrom(inner.Sa, data.NewSimpleDataReadOnly),
		_Ma:     readonly.NewMapFrom(inner.Ma, data.NewReferencedDataReadOnly),
		_Da:     data.NewSimpleDataReadOnly(inner.Da),
		_Age:    readonly.NewPointer(inner.Age),
		_State2: readonly.NewPointer(inner.State2),
		_Data:   readonly.NewList(inner.Data),
	}
}

func (x *PlayerReadOnly) GetName() (_ string) {
	if x == nil || x.inner == nil {
		return
	}
	return x.inner.GetName()
}

func (x *PlayerReadOnly) GetSa() (_ readonly.List[*data.SimpleDataReadOnly]) {
	if x == nil || x.inner == nil {
		return
	}
	return x._Sa
}

func (x *PlayerReadOnly) GetMa() (_ readonly.Map[int32, *data.ReferencedDataReadOnly]) {
	if x == nil || x.inner == nil {
		return
	}
	return x._Ma
}

func (x *PlayerReadOnly) GetDa() (_ *data.SimpleDataReadOnly) {
	if x == nil || x.inner == nil {
		return
	}
	return x._Da
}

func (x *PlayerReadOnly) GetAge() (_ readonly.Pointer[int32]) {
	if x == nil || x.inner == nil {
		return
	}
	return x._Age
}

func (x *PlayerReadOnly) GetState() (_ State) {
	if x == nil || x.inner == nil {
		return
	}
	return x.inner.GetState()
}

func (x *PlayerReadOnly) GetState2() (_ readonly.Pointer[State]) {
	if x == nil || x.inner == nil {
		return
	}
	return x._State2
}

func (x *PlayerReadOnly) GetData() (_ readonly.List[byte]) {
	if x == nil || x.inner == nil {
		return
	}
	return x._Data
}
