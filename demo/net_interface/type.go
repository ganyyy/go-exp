package netinterface

import "fmt"

type IData interface {
	IDataMeta
	Do(string)
}

type IDataMeta interface {
	GetType() DataType
	New() IData
}

type IntData struct {
	Data int
}

// IntData implements IDataTypeMeta
func (i *IntData) GetType() DataType { return DataTypeInt }
func (i *IntData) New() IData        { return &IntData{} }

// IntData implements IDataType
func (i *IntData) Do(s string) {
	fmt.Printf("int data: %d, %s\n", i.Data, s)
}

type StrData struct {
	Data string
}

// StrData implements IDataTypeMeta
func (s *StrData) GetType() DataType { return DataTypeStr }
func (s *StrData) New() IData        { return &StrData{} }

// StrData implements IDataType
func (s *StrData) Do(s2 string) {
	fmt.Printf("str data: %s, %s\n", s.Data, s2)
}

type DataType uint32

const (
	DataTypeInt DataType = iota
	DataTypeStr
)

var (
	dataTypeToMeta = map[DataType]IDataMeta{}
)

func RegisterDataTypeMeta(meta IDataMeta) {
	if _, ok := dataTypeToMeta[meta.GetType()]; ok {
		panic(fmt.Sprintf("data type %d already registered", meta.GetType()))
	}
	dataTypeToMeta[meta.GetType()] = meta
}

func init() {
	RegisterDataTypeMeta(&IntData{})
	RegisterDataTypeMeta(&StrData{})
}
