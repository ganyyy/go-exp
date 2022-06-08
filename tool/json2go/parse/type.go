package parse

type (
	FieldType uint32
)

const (
	TypeInt FieldType = 1 << iota // 这是一个特殊位置, 仅用于区分是否可以使用整形表示
	TypeFloat
	TypeBool
	TypeString
	TypeMap
	TypeObject
	TypeSlice // Slice比较特殊, 因为可能出现多重切片, 这里单独处理一下
	TypeCount = iota
)

func (t FieldType) Check(ft FieldType) bool {
	return t&ft != 0
}

func (t FieldType) Clear(ft FieldType) FieldType {
	return t &^ ft
}

func (t FieldType) NaiveType() FieldType {
	t = t.Clear(TypeSlice)
	t.ClearSliceCount()
	return t.Clear(TypeInt | TypeFloat)
}

func (t FieldType) NumberType() FieldType {
	var ret FieldType
	if t.Check(TypeInt) {
		ret |= TypeInt
	}
	if t.Check(TypeFloat) {
		ret |= TypeFloat
	}
	return ret
}

func (t FieldType) ElemType() FieldType {
	if typ := t.NaiveType(); typ != 0 {
		return typ
	}
	// 再看一下数字类型. 整形和浮点型可能同时存在, 此时以浮点型为准
	if t.Check(TypeFloat) {
		return TypeFloat
	}
	return TypeInt
}

func (t FieldType) FiledType() string {
	var isSlice = t.Check(TypeSlice)
	var typ string
	switch t.ElemType() {
	case TypeInt:
		typ = "int"
	case TypeFloat, TypeFloat | TypeInt:
		typ = "float64"
	case TypeBool:
		typ = "bool"
	case TypeString:
		typ = "string"
	case TypeMap:
		typ = "" // 外围决定对应的类型
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

func (t FieldType) Default() string {
	if t.Check(TypeSlice | TypeObject | TypeMap) {
		return "nil"
	}
	switch t.ElemType() {
	case TypeInt, TypeFloat, TypeInt | TypeFloat:
		return "0"
	case TypeBool:
		return "false"
	case TypeString:
		return "\"\""
	default:
		return "nil"
	}
}

// 以下是针对于Slice的特殊方法

func (t FieldType) SliceCount() int {
	if !t.Check(TypeSlice) {
		return 0
	}
	return int(t >> (TypeCount))
}

func (t *FieldType) SetSlice() {
	if t.Check(TypeSlice) {
		return
	}
	*t |= TypeSlice
	t.AddSlice(1)
}

func (t *FieldType) AddSlice(cnt int) {
	if !t.Check(TypeSlice) || cnt == 0 {
		return
	}
	var total = cnt + t.SliceCount()
	var cur = *t
	cur &^= FieldType(^(uint32(1<<TypeCount) - 1))
	cur |= FieldType(uint32(total << TypeCount))
	*t = cur
}

func (t *FieldType) ClearSliceCount() {
	*t &^= FieldType(^(uint32(1<<TypeCount) - 1))
}

// 针对Map的特殊处理

func (t *FieldType) SetMap() {
	*t = t.Clear(TypeObject) | TypeMap
}
