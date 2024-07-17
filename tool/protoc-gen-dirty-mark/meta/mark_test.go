package meta

import (
	"encoding/json"
	"testing"

	"protoc-gen-dirty-mark/pb"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

type MemoryData struct {
	Mark[*pb.Data]
	Name      string
	Inner     *MemoryInner
	StrMap    *ValueMap[string, string]
	InnerMap  *ReferenceMap[string, *MemoryInner, *pb.Inner]
	StrList   *ValueList[string]
	InnerList *ReferenceList[*MemoryInner, *pb.Inner]
}

func NewMemoryData() *MemoryData {
	var m MemoryData
	m.Inner = NewMemoryInner()
	m.StrMap = NewValueMap[string, string]()
	m.InnerMap = NewReferenceMap[string, *MemoryInner]()
	m.StrList = NewValueList[string]()
	m.InnerList = NewReferenceList[*MemoryInner]()
	return &m
}

// NewValue creates a new MemoryData.
func (*MemoryData) NewValue() IValue[*pb.Data] {
	return NewMemoryData()
}

// FromProto converts from ProtoData.
func (m *MemoryData) FromProto(p *pb.Data) {
	m.SetName(p.GetName())
	m.GetInner().FromProto(p.Inner)
	m.GetStrMap().FromProto(p.StrMap)
	m.GetInnerMap().FromProto(p.InnerMap)

}

// ToProto converts to ProtoData.
func (m *MemoryData) ToProto() *pb.Data {
	var pd pb.Data
	pd.Name = proto.String(m.Name)
	pd.Inner = m.Inner.ToProto()
	pd.StrMap = m.StrMap.ToProto()
	pd.InnerMap = m.InnerMap.ToProto()
	return &pd
}

// ResetDirty resets the dirty mark.
func (m *MemoryData) ResetDirty() {
	m.Mark.ResetDirty()
	m.Inner.ResetDirty()
}

// GetName gets the name.
func (m *MemoryData) GetName() string {
	return m.Name
}

// SetName sets the name.
func (m *MemoryData) SetName(name string) {
	m.Name = name
	m.dirtyName()
}

func (m *MemoryData) dirtyName() {
	m.Dirty(1, m.setProtoName)
}

// setProtoName applies the name.
func (m *MemoryData) setProtoName(p *pb.Data) {
	p.Name = proto.String(m.Name)
}

// GetInner gets the inner.
func (m *MemoryData) GetInner() *MemoryInner {
	if m.Inner == nil {
		m.Inner = NewMemoryInner()
	}
	m.Inner.Dyeing(m.dirtyInner)
	return m.Inner
}

// SetInner sets the inner.
func (m *MemoryData) SetInner(inner *MemoryInner) {
	m.Inner = inner
	m.Dirty(2, m.setProtoInner)
}

// dirtyInner applies the inner.
func (m *MemoryData) dirtyInner() { m.Dirty(2, m.setProtoInner) }

// setProtoInner applies the inner.
func (m *MemoryData) setProtoInner(p *pb.Data) {
	if p.Inner == nil {
		p.Inner = &pb.Inner{}
	}
	m.Inner.DirtyCollect(p.Inner)
}

// GetStrMap gets the strMap.
func (m *MemoryData) GetStrMap() *ValueMap[string, string] {
	m.StrMap.Dyeing(func() { m.Dirty(3, m.setProtoStrMap) })
	return m.StrMap
}

// SetStrMap sets the strMap.
func (m *MemoryData) SetStrMap(strMap *ValueMap[string, string]) {
	m.StrMap = strMap
	m.Dirty(3, m.setProtoStrMap)
}

// setProtoStrMap applies the strMap.
func (m *MemoryData) setProtoStrMap(p *pb.Data) {
	p.StrMap = m.StrMap.DirtyCollect(p.StrMap)
}

// GetInnerMap gets the innerMap.
func (m *MemoryData) GetInnerMap() *ReferenceMap[string, *MemoryInner, *pb.Inner] {
	m.InnerMap.Dyeing(func() { m.Dirty(4, m.setProtoInnerMap) })
	return m.InnerMap
}

// SetInnerMap sets the innerMap.
func (m *MemoryData) SetInnerMap(innerMap *ReferenceMap[string, *MemoryInner, *pb.Inner]) {
	m.InnerMap = innerMap
	m.Dirty(4, m.setProtoInnerMap)
}

// setProtoInnerMap applies the innerMap.
func (m *MemoryData) setProtoInnerMap(p *pb.Data) {
	p.InnerMap = m.InnerMap.DirtyCollect(p.InnerMap)
}

// GetStrList gets the strList.
func (m *MemoryData) GetStrList() *ValueList[string] {
	m.StrList.Dyeing(func() { m.Dirty(5, m.setProtoStrList) })
	return m.StrList
}

// SetStrList sets the strList.
func (m *MemoryData) SetStrList(strList *ValueList[string]) {
	m.StrList = strList
	m.Dirty(5, m.setProtoStrList)
}

// setProtoStrList applies the strList.
func (m *MemoryData) setProtoStrList(p *pb.Data) {
	p.StrList = m.StrList.DirtyCollect(p.StrList)
}

// GetInnerList gets the innerList.
func (m *MemoryData) GetInnerList() *ReferenceList[*MemoryInner, *pb.Inner] {
	m.InnerList.Dyeing(func() { m.Dirty(6, m.setProtoInnerList) })
	return m.InnerList
}

// SetInnerList sets the innerList.
func (m *MemoryData) SetInnerList(innerList *ReferenceList[*MemoryInner, *pb.Inner]) {
	m.InnerList = innerList
	m.Dirty(6, m.setProtoInnerList)
}

// setProtoInnerList applies the innerList.
func (m *MemoryData) setProtoInnerList(p *pb.Data) {
	p.InnerList = m.InnerList.DirtyCollect(p.InnerList)
}

// DirtyProto applies the dirty mark to the target.
func (m *MemoryData) DirtyProto() *pb.Data {
	var p pb.Data
	m.DirtyCollect(&p)
	m.ResetDirty()
	return &p
}

type MemoryInner struct {
	Mark[*pb.Inner]
	Data string
}

func NewMemoryInner() *MemoryInner {
	var m MemoryInner
	return &m
}

// NewValue creates a new MemoryInner.
func (*MemoryInner) NewValue() IValue[*pb.Inner] {
	return NewMemoryInner()
}

// From
func (m *MemoryInner) FromProto(p *pb.Inner) {
	m.dyed()
	m.Data = p.Data
}

// To
func (m *MemoryInner) ToProto() *pb.Inner {
	var proto pb.Inner
	proto.Data = m.Data
	return &proto
}

// GetData gets the data.
func (m *MemoryInner) GetData() string {
	return m.Data
}

// SetData sets the data.
func (m *MemoryInner) SetData(data string) {
	m.Data = data
	m.Dirty(1, m.setProtoData)
}

// setProtoData applies the data.
func (m *MemoryInner) setProtoData(p *pb.Inner) {
	p.Data = m.Data
}

func TestMark(t *testing.T) {
	var m = NewMemoryData()
	m.SetName("test")
	m.SetName("test2")
	m.SetName("test3")

	inner := m.GetInner()
	inner.SetData("inner")
	inner.SetData("inner2")
	inner.SetData("inner3")

	strMap := m.GetStrMap()
	strMap.Set("key", "value")
	strMap.Set("key2", "value2")
	strMap.Del("key")

	innerMap := m.GetInnerMap()
	inner1 := NewMemoryInner()
	inner1.SetData("inner1")
	inner2 := NewMemoryInner()
	inner2.SetData("inner2")
	innerMap.Set("key1", inner1)
	innerMap.Set("key2", inner2)
	innerMap.Del("key2")

	inner1 = innerMap.Get("key1")
	inner1.SetData("inner1-1")

	l1 := m.GetStrList()
	l1.Add("str1", "str2", "str3")

	l2 := m.GetInnerList()
	lInner1 := NewMemoryInner()
	lInner1.SetData("lInner1")
	lInner2 := NewMemoryInner()
	lInner2.SetData("lInner2")
	l2.Add(lInner1, lInner2)

	var p pb.Data
	m.DirtyCollect(&p)
	m.ResetDirty()
	bs, _ := json.Marshal(&p)
	t.Log(string(bs))

	m.SetName("test4")
	l1 = m.GetStrList()
	l1.Add("str4")
	m.GetInner().SetData("inner4")
	m1 := m.GetInnerMap()
	m1.Set("key3", inner1)
	p = pb.Data{}
	m.DirtyCollect(&p)
	m.ResetDirty()
	bs, _ = json.Marshal(&p)
	t.Log(string(bs))

	// _ = proto.Unmarshal(bs, &p)

	// input := protoiface.UnmarshalInput{
	// 	Message:  p.ProtoReflect(),
	// 	Buf:      bs,
	// 	Resolver: protoregistry.GlobalTypes,
	// 	Depth:    protowire.DefaultRecursionLimit,
	// }

}

func TestBitSet(t *testing.T) {
	var b = NewBitsetMark(100)
	b.Dirty(0)
	b.Dirty(1)
	b.Dirty(2)
	b.Dirty(3)
	b.Dirty(50)
	b.Dirty(99)
	b.Dirty(100)
	b.Dirty(101)

	for i := range b.AllBits() {
		t.Log(i)
	}
}

func TestProtoWire(t *testing.T) {
	var p pb.Data
	p.Name = proto.String("name")
	p.Inner = &pb.Inner{
		Data: "inner",
		Age:  1,
	}
	p.StrMap = map[string]string{
		"key1": "value1",
		"key2": "value26",
	}
	p.InnerMap = map[string]*pb.Inner{
		"key1": {Data: "inner133", Age: 123456},
		"key2": {Data: "inner2444", Age: 1},
	}
	p.StrList = []string{"str1", "str25"}
	p.InnerList = []*pb.Inner{
		{Data: "inner11"},
		{Data: "inner2222"},
	}

	// p.Reset()
	// p.Inner = &pb.Inner{
	// 	Data: "inner",
	// 	Age:  1,
	// }

	bs, _ := proto.Marshal(&p)
	_ = proto.Unmarshal(bs, &p)

	pr := p.ProtoReflect()

	fields := pr.Descriptor().Fields()

	_ = fields

	var fieldsData = make(map[protowire.Number][]byte)

	for len(bs) > 0 {
		num, typ, headLength := protowire.ConsumeTag(bs)
		if headLength < 0 {
			t.Fatal(headLength)
		}
		bodyLength := protowire.ConsumeFieldValue(num, typ, bs[headLength:])
		if bodyLength < 0 {
			t.Fatal(bodyLength)
		}
		total := headLength + bodyLength
		t.Log(num, typ, headLength, bodyLength)
		fieldsData[num] = append(fieldsData[num], bs[:total]...)
		bs = bs[total:]
	}

	var bs2 []byte
	bs2 = append(bs2, fieldsData[2]...)
	bs2 = append(bs2, fieldsData[32345]...)
	bs2 = append(bs2, fieldsData[62231]...)

	_ = proto.Unmarshal(bs2, &p)
	t.Logf("%+v", p.String())
}
