package test

import "func-exportor/types"

// 测试用的导出符号
const TestConstant string = "test"

var (
	TestVariable   int            = 42
	Variable2      string         = "hello"
	Variable3      float32        = 3.14
	NotExportedVar                = "not exported"
	OutMap         map[string]int = map[string]int{}
	OutSlice       []string       = []string{"one", "two", "three"}
	OutBox         Box[int]       = Box[int]{Item: 100}
)

type Box[T any] struct {
	Item T
}

type TestStruct struct {
	Name string
}

type TestInterface interface {
	Method() string
}

func TestFunction() string {
	return "test"
}

func TestFunction2(a int, b string) (int, string) {
	return a + 1, b + "!"
}

func (t *TestStruct) ExportedMethod() string {
	return t.Name
}

func ExportType(name string) types.ExportedSymbol {
	return types.ExportedSymbol{
		Name: name,
		Type: "TestStruct",
	}
}

func NewTestStruct(name string) *TestStruct {
	return &TestStruct{Name: name}
}

// 非导出符号
const internal = "internal"

var privateVar = 10

type internalType struct {
	data string
}

func internalFunc() {
	// do nothing
}

func TestStructSetName(ts *TestStruct, name string) {
	ts.Name = name
}
