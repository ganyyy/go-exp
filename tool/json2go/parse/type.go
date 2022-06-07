package parse

type (
	FiledType uint32
)

const (
	TypeInt FiledType = 1 << iota // 这是一个特殊位置, 仅用于区分是否可以使用整形表示
	TypeFloat
	TypeBool
	TypeString
	TypeObject
	TypeSlice // Slice比较特殊, 因为可能出现多重切片, 这里单独处理一下
	TypeCount = iota
)

func (t FiledType) Check(ft FiledType) bool {
	return t&ft != 0
}

func (t FiledType) Clear(ft FiledType) FiledType {
	return t &^ ft
}

func (t FiledType) NaiveType() FiledType {
	t = t.Clear(TypeSlice)
	t.ClearSliceCount()
	return t.Clear(TypeInt | TypeFloat)
}

func (t FiledType) NumberType() FiledType {
	var ret FiledType
	if t.Check(TypeInt) {
		ret |= TypeInt
	}
	if t.Check(TypeFloat) {
		ret |= TypeFloat
	}
	return ret
}

func (t FiledType) ElemType() FiledType {
	if typ := t.NaiveType(); typ != 0 {
		return typ
	}
	// 再看一下数字类型. 整形和浮点型可能同时存在, 此时以浮点型为准
	if t.Check(TypeFloat) {
		return TypeFloat
	}
	return TypeInt
}

func (t FiledType) FiledType() string {
	var isSlice = t.Check(TypeSlice)
	var typ string
	switch t.ElemType() {
	case TypeInt:
		typ = "int"
	case TypeFloat:
		typ = "float64"
	case TypeBool:
		typ = "bool"
	case TypeString:
		typ = "string"
	case TypeObject:
		// 结构体类型的时机名称需要在外部添加
		typ = "*"
	}
	var pre string
	if isSlice {
		// 感谢类型前置...
		for i := 0; i < t.SliceCount(); i++ {
			pre += "[]"
		}
	}
	return pre + typ
}

func (t FiledType) Default() string {
	if t.Check(TypeSlice) {
		return "nil"
	}
	switch t.ElemType() {
	case TypeInt, TypeFloat:
		return "0"
	case TypeBool:
		return "false"
	case TypeString:
		return "\"\""
	case TypeObject:
		fallthrough
	default:
		return "nil"
	}
}

// 以下是针对于Slice的特殊方法

func (t FiledType) SliceCount() int {
	if !t.Check(TypeSlice) {
		return 0
	}
	return int(t >> (TypeCount))
}

func (t *FiledType) SetSlice() {
	if t.Check(TypeSlice) {
		return
	}
	*t |= TypeSlice
	t.AddSlice(1)
}

func (t *FiledType) AddSlice(cnt int) {
	if !t.Check(TypeSlice) || cnt == 0 {
		return
	}
	var total = cnt + t.SliceCount()
	var cur = *t
	cur &^= FiledType(^(uint32(1<<TypeCount) - 1))
	cur |= FiledType(uint32(total << TypeCount))
	*t = cur
}

func (t *FiledType) ClearSliceCount() {
	*t &^= FiledType(^(uint32(1<<TypeCount) - 1))
}
