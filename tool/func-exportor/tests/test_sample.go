package test

// 测试用的导出符号
const TestConstant = "test"

var TestVariable = 42

type TestStruct struct {
	Name string
}

type TestInterface interface {
	Method() string
}

func TestFunction() string {
	return "test"
}

func (t *TestStruct) ExportedMethod() string {
	return t.Name
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
