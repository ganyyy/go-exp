package data

import (
	meta "protoc-gen-dirty-mark/meta"
	pb "protoc-gen-dirty-mark/pb"
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

var _DataApplyDirtyTable = []func(*Data, *pb.Data){
	DataFieldIndexName:      (*Data).applyDirtyName,
	DataFieldIndexInner:     (*Data).applyDirtyInner,
	DataFieldIndexStrMap:    (*Data).applyDirtyStrMap,
	DataFieldIndexInnerMap:  (*Data).applyDirtyInnerMap,
	DataFieldIndexStrList:   (*Data).applyDirtyStrList,
	DataFieldIndexInnerList: (*Data).applyDirtyInnerList,
}

type Data struct {
	mark       *meta.BitsetMark
	_Name      string
	_Inner     *Inner
	_StrMap    *meta.ValueMap[string, string]
	_InnerMap  *meta.ReferenceMap[string, *Inner, *pb.Inner]
	_StrList   *meta.ValueList[string]
	_InnerList *meta.ReferenceList[*Inner, *pb.Inner]
}

func NewData() *Data {
	var m Data
	m.mark = meta.NewBitsetMark(DataFieldMax)
	m._Inner = NewInner()
	m._StrMap = meta.NewValueMap[string, string]()
	m._InnerMap = meta.NewReferenceMap[string, *Inner, *pb.Inner]()
	m._StrList = meta.NewValueList[string]()
	m._InnerList = meta.NewReferenceList[*Inner, *pb.Inner]()
	return &m
}

// NewValue creates a new Data.
func (*Data) NewValue() meta.IValue[*pb.Data] {
	return NewData()
}

// GetName gets the Name.
func (m *Data) GetName() string {
	return m._Name
}

// SetName sets the Name.
func (m *Data) SetName(v string) {
	m._Name = v
	meta.MarkHelper(m.mark, DataFieldIndexName)
}

func (m *Data) applyDirtyName(p *pb.Data) {
	p.Name = meta.Pointer(m.GetName())
}

// GetInner gets the Inner.
func (m *Data) GetInner() *Inner {
	if m._Inner == nil {
		m._Inner = NewInner()
	}
	meta.SetMarkHelper(m._Inner.mark, m.mark, DataFieldIndexInner)
	return m._Inner
}

// SetInner sets the Inner.
func (m *Data) SetInner(v *Inner) {
	if m._Inner == v {
		return
	}
	if m._Inner != nil {
		meta.SetMarkHelper(m._Inner.mark, nil, 0)
	}
	m._Inner = v
	if v != nil {
		meta.SetMarkHelper(m._Inner.mark, m.mark, DataFieldIndexInner)
	}
	meta.MarkHelper(m.mark, DataFieldIndexInner)
}

func (m *Data) applyDirtyInner(p *pb.Data) {
	if p.Inner == nil {
		p.Inner = &pb.Inner{}
	}
	m.GetInner().DirtyCollect(p.Inner)
}

// GetStrMap gets the StrMap.
func (m *Data) GetStrMap() *meta.ValueMap[string, string] {
	if m._StrMap == nil {
		m._StrMap = meta.NewValueMap[string, string]()
	}
	meta.SetMarkHelper(m._StrMap, m.mark, DataFieldIndexStrMap)
	return m._StrMap
}

// SetStrMap sets the StrMap.
func (m *Data) SetStrMap(v *meta.ValueMap[string, string]) {
	if m._StrMap == v {
		return
	}
	if m._StrMap != nil {
		meta.SetMarkHelper(m._StrMap, nil, 0)
	}
	m._StrMap = v
	if v != nil {
		meta.SetMarkHelper(m._StrMap, m.mark, DataFieldIndexStrMap)
	}
	meta.MarkHelper(m.mark, DataFieldIndexStrMap)
}

func (m *Data) applyDirtyStrMap(p *pb.Data) {
	p.StrMap = m.GetStrMap().DirtyCollect(p.StrMap)
}

// GetInnerMap gets the InnerMap.
func (m *Data) GetInnerMap() *meta.ReferenceMap[string, *Inner, *pb.Inner] {
	if m._InnerMap == nil {
		m._InnerMap = meta.NewReferenceMap[string, *Inner, *pb.Inner]()
	}
	meta.SetMarkHelper(m._InnerMap, m.mark, DataFieldIndexInnerMap)
	return m._InnerMap
}

// SetInnerMap sets the InnerMap.
func (m *Data) SetInnerMap(v *meta.ReferenceMap[string, *Inner, *pb.Inner]) {
	if m._InnerMap == v {
		return
	}
	if m._InnerMap != nil {
		meta.SetMarkHelper(m._InnerMap, nil, 0)
	}
	m._InnerMap = v
	if v != nil {
		meta.SetMarkHelper(m._InnerMap, m.mark, DataFieldIndexInnerMap)
	}
	meta.MarkHelper(m.mark, DataFieldIndexInnerMap)
}

func (m *Data) applyDirtyInnerMap(p *pb.Data) {
	p.InnerMap = m.GetInnerMap().DirtyCollect(p.InnerMap)
}

// GetStrList gets the StrList.
func (m *Data) GetStrList() *meta.ValueList[string] {
	if m._StrList == nil {
		m._StrList = meta.NewValueList[string]()
	}
	meta.SetMarkHelper(m._StrList, m.mark, DataFieldIndexStrList)
	return m._StrList
}

// SetStrList sets the StrList.
func (m *Data) SetStrList(v *meta.ValueList[string]) {
	if m._StrList == v {
		return
	}
	if m._StrList != nil {
		meta.SetMarkHelper(m._StrList, nil, 0)
	}
	m._StrList = v
	if v != nil {
		meta.SetMarkHelper(m._StrList, m.mark, DataFieldIndexStrList)
	}
	meta.MarkHelper(m.mark, DataFieldIndexStrList)
}

func (m *Data) applyDirtyStrList(p *pb.Data) {
	p.StrList = m.GetStrList().DirtyCollect(p.StrList)
}

// GetInnerList gets the InnerList.
func (m *Data) GetInnerList() *meta.ReferenceList[*Inner, *pb.Inner] {
	if m._InnerList == nil {
		m._InnerList = meta.NewReferenceList[*Inner, *pb.Inner]()
	}
	meta.SetMarkHelper(m._InnerList, m.mark, DataFieldIndexInnerList)
	return m._InnerList
}

// SetInnerList sets the InnerList.
func (m *Data) SetInnerList(v *meta.ReferenceList[*Inner, *pb.Inner]) {
	if m._InnerList == v {
		return
	}
	if m._InnerList != nil {
		meta.SetMarkHelper(m._InnerList, nil, 0)
	}
	m._InnerList = v
	if v != nil {
		meta.SetMarkHelper(m._InnerList, m.mark, DataFieldIndexInnerList)
	}
	meta.MarkHelper(m.mark, DataFieldIndexInnerList)
}

func (m *Data) applyDirtyInnerList(p *pb.Data) {
	p.InnerList = m.GetInnerList().DirtyCollect(p.InnerList)
}

// FromProto sets the value from the target.
func (m *Data) FromProto(p *pb.Data) {
	m.SetName(p.GetName())
	m.GetInner().FromProto(p.GetInner())
	m.GetStrMap().FromProto(p.GetStrMap())
	m.GetInnerMap().FromProto(p.GetInnerMap())
	m.GetStrList().FromProto(p.GetStrList())
	m.GetInnerList().FromProto(p.GetInnerList())
}

// ToProto gets the target from the value.
func (m *Data) ToProto() *pb.Data {
	var p pb.Data
	p.Name = meta.Pointer(m.GetName())
	p.Inner = m.GetInner().ToProto()
	p.StrMap = m.GetStrMap().ToProto()
	p.InnerMap = m.GetInnerMap().ToProto()
	p.StrList = m.GetStrList().ToProto()
	p.InnerList = m.GetInnerList().ToProto()
	return &p
}

// resetDirty resets the dirty mark.
func (m *Data) resetDirty() {
	meta.ResetHelper(m.mark)
	m.GetInner().resetDirty()
	meta.ResetHelper(m.GetStrMap())
	meta.ResetHelper(m.GetInnerMap())
	meta.ResetHelper(m.GetStrList())
	meta.ResetHelper(m.GetInnerList())
}

// DirtyProto returns proto apply the dirty mark.
func (m *Data) DirtyProto() *pb.Data {
	var p pb.Data
	m.DirtyCollect(&p)
	return &p
}

// DirtyCollect applies the dirty mark to the target.
func (m *Data) DirtyCollect(target *pb.Data) {
	for dirtyIdx := range meta.DirtyBitsHelper(m.mark) {
		_DataApplyDirtyTable[dirtyIdx](m, target)
	}
	m.resetDirty()
}

// GetMark gets the mark.
func (m *Data) GetMark() meta.IMark {
	return m.mark
}

const (
	_Data2FieldIndex = iota - 1
	Data2FieldIndexId
	Data2FieldIndexData
	Data2FieldMax
)

var _Data2ApplyDirtyTable = []func(*Data2, *pb.Data2){
	Data2FieldIndexId:   (*Data2).applyDirtyId,
	Data2FieldIndexData: (*Data2).applyDirtyData,
}

type Data2 struct {
	mark  *meta.BitsetMark
	_Id   int32
	_Data *meta.ValueList[byte]
}

func NewData2() *Data2 {
	var m Data2
	m.mark = meta.NewBitsetMark(Data2FieldMax)
	m._Data = meta.NewValueList[byte]()
	return &m
}

// NewValue creates a new Data2.
func (*Data2) NewValue() meta.IValue[*pb.Data2] {
	return NewData2()
}

// GetId gets the Id.
func (m *Data2) GetId() int32 {
	return m._Id
}

// SetId sets the Id.
func (m *Data2) SetId(v int32) {
	m._Id = v
	meta.MarkHelper(m.mark, Data2FieldIndexId)
}

func (m *Data2) applyDirtyId(p *pb.Data2) {
	p.Id = meta.Pointer(m.GetId())
}

// GetData gets the Data.
func (m *Data2) GetData() *meta.ValueList[byte] {
	if m._Data == nil {
		m._Data = meta.NewValueList[byte]()
	}
	meta.SetMarkHelper(m._Data, m.mark, Data2FieldIndexData)
	return m._Data
}

// SetData sets the Data.
func (m *Data2) SetData(v *meta.ValueList[byte]) {
	if m._Data == v {
		return
	}
	if m._Data != nil {
		meta.SetMarkHelper(m._Data, nil, 0)
	}
	m._Data = v
	if v != nil {
		meta.SetMarkHelper(m._Data, m.mark, Data2FieldIndexData)
	}
	meta.MarkHelper(m.mark, Data2FieldIndexData)
}

func (m *Data2) applyDirtyData(p *pb.Data2) {
	p.Data = m.GetData().DirtyCollect(p.Data)
}

// FromProto sets the value from the target.
func (m *Data2) FromProto(p *pb.Data2) {
	m.SetId(p.GetId())
	m.GetData().FromProto(p.GetData())
}

// ToProto gets the target from the value.
func (m *Data2) ToProto() *pb.Data2 {
	var p pb.Data2
	p.Id = meta.Pointer(m.GetId())
	p.Data = m.GetData().ToProto()
	return &p
}

// resetDirty resets the dirty mark.
func (m *Data2) resetDirty() {
	meta.ResetHelper(m.mark)
	meta.ResetHelper(m.GetData())
}

// DirtyProto returns proto apply the dirty mark.
func (m *Data2) DirtyProto() *pb.Data2 {
	var p pb.Data2
	m.DirtyCollect(&p)
	return &p
}

// DirtyCollect applies the dirty mark to the target.
func (m *Data2) DirtyCollect(target *pb.Data2) {
	for dirtyIdx := range meta.DirtyBitsHelper(m.mark) {
		_Data2ApplyDirtyTable[dirtyIdx](m, target)
	}
	m.resetDirty()
}

// GetMark gets the mark.
func (m *Data2) GetMark() meta.IMark {
	return m.mark
}

const (
	_InnerFieldIndex = iota - 1
	InnerFieldIndexData
	InnerFieldIndexAge
	InnerFieldMax
)

var _InnerApplyDirtyTable = []func(*Inner, *pb.Inner){
	InnerFieldIndexData: (*Inner).applyDirtyData,
	InnerFieldIndexAge:  (*Inner).applyDirtyAge,
}

type Inner struct {
	mark  *meta.BitsetMark
	_Data string
	_Age  int32
}

func NewInner() *Inner {
	var m Inner
	m.mark = meta.NewBitsetMark(InnerFieldMax)
	return &m
}

// NewValue creates a new Inner.
func (*Inner) NewValue() meta.IValue[*pb.Inner] {
	return NewInner()
}

// GetData gets the Data.
func (m *Inner) GetData() string {
	return m._Data
}

// SetData sets the Data.
func (m *Inner) SetData(v string) {
	m._Data = v
	meta.MarkHelper(m.mark, InnerFieldIndexData)
}

func (m *Inner) applyDirtyData(p *pb.Inner) {
	p.Data = m.GetData()
}

// GetAge gets the Age.
func (m *Inner) GetAge() int32 {
	return m._Age
}

// SetAge sets the Age.
func (m *Inner) SetAge(v int32) {
	m._Age = v
	meta.MarkHelper(m.mark, InnerFieldIndexAge)
}

func (m *Inner) applyDirtyAge(p *pb.Inner) {
	p.Age = m.GetAge()
}

// FromProto sets the value from the target.
func (m *Inner) FromProto(p *pb.Inner) {
	m.SetData(p.GetData())
	m.SetAge(p.GetAge())
}

// ToProto gets the target from the value.
func (m *Inner) ToProto() *pb.Inner {
	var p pb.Inner
	p.Data = m.GetData()
	p.Age = m.GetAge()
	return &p
}

// resetDirty resets the dirty mark.
func (m *Inner) resetDirty() {
	meta.ResetHelper(m.mark)
}

// DirtyProto returns proto apply the dirty mark.
func (m *Inner) DirtyProto() *pb.Inner {
	var p pb.Inner
	m.DirtyCollect(&p)
	return &p
}

// DirtyCollect applies the dirty mark to the target.
func (m *Inner) DirtyCollect(target *pb.Inner) {
	for dirtyIdx := range meta.DirtyBitsHelper(m.mark) {
		_InnerApplyDirtyTable[dirtyIdx](m, target)
	}
	m.resetDirty()
}

// GetMark gets the mark.
func (m *Inner) GetMark() meta.IMark {
	return m.mark
}
