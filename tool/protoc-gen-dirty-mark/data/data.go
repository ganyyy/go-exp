package data

import (
	meta1 "protoc-gen-dirty-mark/meta"
	pb123 "protoc-gen-dirty-mark/pb"
)

const (
	_DataFieldIndex = iota - 1
	DataFieldIndexName
	DataFieldIndexInner
	DataFieldIndexStrMap
	DataFieldIndexInnerMap
	DataFieldIndexStrList
	DataFieldIndexInnerList
	DataFieldMax
)

var _DataApplyDirtyTable = []func(*Data, *pb123.Data){
	DataFieldIndexName:      (*Data).applyDirtyName,
	DataFieldIndexInner:     (*Data).applyDirtyInner,
	DataFieldIndexStrMap:    (*Data).applyDirtyStrMap,
	DataFieldIndexInnerMap:  (*Data).applyDirtyInnerMap,
	DataFieldIndexStrList:   (*Data).applyDirtyStrList,
	DataFieldIndexInnerList: (*Data).applyDirtyInnerList,
}

type Data struct {
	mark       *meta1.BitsetMark
	_Name      string
	_Inner     *Inner
	_StrMap    *meta1.ValueMap[string, string]
	_InnerMap  *meta1.ReferenceMap[string, *Inner, *pb123.Inner]
	_StrList   *meta1.ValueList[string]
	_InnerList *meta1.ReferenceList[*Inner, *pb123.Inner]
}

func NewData() *Data {
	var m Data
	m.mark = meta1.NewBitsetMark(DataFieldMax)
	m._Inner = NewInner()
	m._StrMap = meta1.NewValueMap[string, string]()
	m._InnerMap = meta1.NewReferenceMap[string, *Inner, *pb123.Inner]()
	m._StrList = meta1.NewValueList[string]()
	m._InnerList = meta1.NewReferenceList[*Inner, *pb123.Inner]()
	return &m
}

// NewValue creates a new Data.
func (*Data) NewValue() meta1.IValue[*pb123.Data] {
	return NewData()
}

// GetName gets the Name.
func (m *Data) GetName() string {
	return m._Name
}

// SetName sets the Name.
func (m *Data) SetName(v string) {
	m._Name = v
	m.dirtyName()
}

func (m *Data) dirtyName() { m.mark.Dirty(DataFieldIndexName) }

func (m *Data) applyDirtyName(p *pb123.Data) {
	p.Name = m.GetName()
}

// GetInner gets the Inner.
func (m *Data) GetInner() *Inner {
	if m._Inner == nil {
		m._Inner = NewInner()
	}
	m._Inner.Dyeing(m.dirtyInner)
	return m._Inner
}

// SetInner sets the Inner.
func (m *Data) SetInner(v *Inner) {
	m._Inner = v
	if v != nil {
		v.Dyeing(m.dirtyInner)
	}
	m.dirtyInner()
}

func (m *Data) dirtyInner() { m.mark.Dirty(DataFieldIndexInner) }

func (m *Data) applyDirtyInner(p *pb123.Data) {
	if p.Inner == nil {
		p.Inner = &pb123.Inner{}
	}
	m.GetInner().DirtyCollect(p.Inner)
}

// GetStrMap gets the StrMap.
func (m *Data) GetStrMap() *meta1.ValueMap[string, string] {
	if m._StrMap == nil {
		m._StrMap = meta1.NewValueMap[string, string]()
	}
	m._StrMap.Dyeing(m.dirtyStrMap)
	return m._StrMap
}

// SetStrMap sets the StrMap.
func (m *Data) SetStrMap(v *meta1.ValueMap[string, string]) {
	m._StrMap = v
	if v != nil {
		v.Dyeing(m.dirtyStrMap)
	}
	m.dirtyStrMap()
}

func (m *Data) dirtyStrMap() { m.mark.Dirty(DataFieldIndexStrMap) }

func (m *Data) applyDirtyStrMap(p *pb123.Data) {
	p.StrMap = m.GetStrMap().DirtyCollect(p.StrMap)
}

// GetInnerMap gets the InnerMap.
func (m *Data) GetInnerMap() *meta1.ReferenceMap[string, *Inner, *pb123.Inner] {
	if m._InnerMap == nil {
		m._InnerMap = meta1.NewReferenceMap[string, *Inner, *pb123.Inner]()
	}
	m._InnerMap.Dyeing(m.dirtyInnerMap)
	return m._InnerMap
}

// SetInnerMap sets the InnerMap.
func (m *Data) SetInnerMap(v *meta1.ReferenceMap[string, *Inner, *pb123.Inner]) {
	m._InnerMap = v
	if v != nil {
		v.Dyeing(m.dirtyInnerMap)
	}
	m.dirtyInnerMap()
}

func (m *Data) dirtyInnerMap() { m.mark.Dirty(DataFieldIndexInnerMap) }

func (m *Data) applyDirtyInnerMap(p *pb123.Data) {
	p.InnerMap = m.GetInnerMap().DirtyCollect(p.InnerMap)
}

// GetStrList gets the StrList.
func (m *Data) GetStrList() *meta1.ValueList[string] {
	if m._StrList == nil {
		m._StrList = meta1.NewValueList[string]()
	}
	m._StrList.Dyeing(m.dirtyStrList)
	return m._StrList
}

// SetStrList sets the StrList.
func (m *Data) SetStrList(v *meta1.ValueList[string]) {
	m._StrList = v
	if v != nil {
		v.Dyeing(m.dirtyStrList)
	}
	m.dirtyStrList()
}

func (m *Data) dirtyStrList() { m.mark.Dirty(DataFieldIndexStrList) }

func (m *Data) applyDirtyStrList(p *pb123.Data) {
	p.StrList = m.GetStrList().DirtyCollect(p.StrList)
}

// GetInnerList gets the InnerList.
func (m *Data) GetInnerList() *meta1.ReferenceList[*Inner, *pb123.Inner] {
	if m._InnerList == nil {
		m._InnerList = meta1.NewReferenceList[*Inner, *pb123.Inner]()
	}
	m._InnerList.Dyeing(m.dirtyInnerList)
	return m._InnerList
}

// SetInnerList sets the InnerList.
func (m *Data) SetInnerList(v *meta1.ReferenceList[*Inner, *pb123.Inner]) {
	m._InnerList = v
	if v != nil {
		v.Dyeing(m.dirtyInnerList)
	}
	m.dirtyInnerList()
}

func (m *Data) dirtyInnerList() { m.mark.Dirty(DataFieldIndexInnerList) }

func (m *Data) applyDirtyInnerList(p *pb123.Data) {
	p.InnerList = m.GetInnerList().DirtyCollect(p.InnerList)
}

// FromProto sets the value from the target.
func (m *Data) FromProto(p *pb123.Data) {
	m.SetName(p.GetName())
	m.GetInner().FromProto(p.GetInner())
	m.GetStrMap().FromProto(p.GetStrMap())
	m.GetInnerMap().FromProto(p.GetInnerMap())
	m.GetStrList().FromProto(p.GetStrList())
	m.GetInnerList().FromProto(p.GetInnerList())
}

// ToProto gets the target from the value.
func (m *Data) ToProto() *pb123.Data {
	var p pb123.Data
	p.Name = m.GetName()
	p.Inner = m.GetInner().ToProto()
	p.StrMap = m.GetStrMap().ToProto()
	p.InnerMap = m.GetInnerMap().ToProto()
	p.StrList = m.GetStrList().ToProto()
	p.InnerList = m.GetInnerList().ToProto()
	return &p
}

// ResetDirty resets the dirty mark.
func (m *Data) ResetDirty() {
	m.mark.Reset()
	m.GetInner().ResetDirty()
}

// DirtyProto returns proto apply the dirty mark.
func (m *Data) DirtyProto() *pb123.Data {
	var p pb123.Data
	m.DirtyCollect(&p)
	return &p
}

// Dyeing set the dyeing function.
func (m *Data) Dyeing(d func()) {
	m.mark.Dyeing(d)
}

// DirtyCollect applies the dirty mark to the target.
func (m *Data) DirtyCollect(target *pb123.Data) {
	for dirtyIdx := range m.mark.AllBits() {
		_DataApplyDirtyTable[dirtyIdx](m, target)
	}
	m.ResetDirty()
}

const (
	_InnerFieldIndex = iota - 1
	InnerFieldIndexData
	InnerFieldMax
)

var _InnerApplyDirtyTable = []func(*Inner, *pb123.Inner){
	InnerFieldIndexData: (*Inner).applyDirtyData,
}

type Inner struct {
	mark  *meta1.BitsetMark
	_Data string
}

func NewInner() *Inner {
	var m Inner
	m.mark = meta1.NewBitsetMark(InnerFieldMax)
	return &m
}

// NewValue creates a new Inner.
func (*Inner) NewValue() meta1.IValue[*pb123.Inner] {
	return NewInner()
}

// GetData gets the Data.
func (m *Inner) GetData() string {
	return m._Data
}

// SetData sets the Data.
func (m *Inner) SetData(v string) {
	m._Data = v
	m.dirtyData()
}

func (m *Inner) dirtyData() { m.mark.Dirty(InnerFieldIndexData) }

func (m *Inner) applyDirtyData(p *pb123.Inner) {
	p.Data = m.GetData()
}

// FromProto sets the value from the target.
func (m *Inner) FromProto(p *pb123.Inner) {
	m.SetData(p.GetData())
}

// ToProto gets the target from the value.
func (m *Inner) ToProto() *pb123.Inner {
	var p pb123.Inner
	p.Data = m.GetData()
	return &p
}

// ResetDirty resets the dirty mark.
func (m *Inner) ResetDirty() {
	m.mark.Reset()
}

// DirtyProto returns proto apply the dirty mark.
func (m *Inner) DirtyProto() *pb123.Inner {
	var p pb123.Inner
	m.DirtyCollect(&p)
	return &p
}

// Dyeing set the dyeing function.
func (m *Inner) Dyeing(d func()) {
	m.mark.Dyeing(d)
}

// DirtyCollect applies the dirty mark to the target.
func (m *Inner) DirtyCollect(target *pb123.Inner) {
	for dirtyIdx := range m.mark.AllBits() {
		_InnerApplyDirtyTable[dirtyIdx](m, target)
	}
	m.ResetDirty()
}
