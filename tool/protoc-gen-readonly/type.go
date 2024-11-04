package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Kind uint8

const (
	KindSimple      Kind = iota // 基础类型, 包括枚举
	KindPointer                 // 基础类型指针, 包括枚举. 用于 optional 字段
	KindMessage                 // 自定义的引用类型. 可能处于不同的包
	KindList                    // 值是简单类型的列表
	KindMessageList             // 值是自定义类型的列表
	KindMap                     // 值是简单类型的 map
	KindMessageMap              // 值是自定义类型的 map
)

type Field struct {
	Kind      Kind
	Type      string // [pkg.]type
	OrgType   string // 原始类型
	Pkg       string
	OuterName string
	InnerName string
}

func (f *Field) FieldString(readOnlyPkg string) string {
	if f.Kind == KindSimple {
		return ""
	}
	return genField(f.InnerName, f.Type)
}

// TypeString returns the string representation of the field
func TypeString(f *Field, readOnlyPkg string) string {
	switch f.Kind {
	case KindSimple, KindMessage:
		return f.Type
	case KindList, KindMessageList:
		return readOnlyPkg + ".List[" + f.Type + "]"
	case KindMap, KindMessageMap:
		return readOnlyPkg + ".Map[" + f.Type + "]"
	case KindPointer:
		return readOnlyPkg + ".Pointer[" + f.Type + "]"
	default:
		panic("unknown kind")
	}
}

func InitFieldString(f *Field, readOnlyPkg string) string {
	formatPkg := func(pkg string) string {
		if pkg == "" {
			return ""
		}
		return pkg + "."
	}
	switch f.Kind {
	case KindSimple:
		return ""
	case KindPointer:
		return fmt.Sprintf("%s: %sNewPointer(inner.%s),",
			f.InnerName, formatPkg(readOnlyPkg), f.OuterName)
	case KindMessage:
		return fmt.Sprintf("%s: %sNew%s(inner.%s),",
			f.InnerName, formatPkg(f.Pkg), readOnlyType(f.OrgType), f.OuterName)
	case KindList, KindMap:
		typ := "NewList"
		if f.Kind == KindMap {
			typ = "NewMap"
		}
		return fmt.Sprintf("%s: %s%s(inner.%s),",
			f.InnerName, formatPkg(readOnlyPkg), typ, f.OuterName)
	case KindMessageList, KindMessageMap:
		typ := "NewListFrom"
		if f.Kind == KindMessageMap {
			typ = "NewMapFrom"
		}
		return fmt.Sprintf("%s: %s%s(inner.%s, %sNew%s),",
			f.InnerName, formatPkg(readOnlyPkg), typ, f.OuterName, formatPkg(f.Pkg), readOnlyType(f.OrgType))
	default:
		panic("unknown kind")
	}
}

type Message struct {
	Name    string
	PkgName string
	Fields  []*Field // 所有字段
	// NOT SUPPORT ONEOF AND EXTENSION
}

type File struct {
	Name     string
	Messages []*Message // 所有消息
}

const (
	ReadOnlyMsgSuffix = "ReadOnly"
)

// AddMessage
func (f *File) AddMessage(name string) *Message {
	msg := &Message{Name: readOnlyType(name)}
	f.Messages = append(f.Messages, msg)
	return msg
}

type GenCfg struct {
	ReadOnlyPkg string
	*protogen.GeneratedFile
}

func (m *Message) AddField(field *protogen.Field, gen GenCfg) {
	// 返回字段的类型和是否是指针
	var getType = func(field *protogen.Field) (string, bool) {
		var goType string
		var pointer = field.Desc.HasPresence() // optional fields are pointers
		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			goType = "bool"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			goType = "int32"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			goType = "uint32"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			goType = "int64"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			goType = "uint64"
		case protoreflect.FloatKind:
			goType = "float32"
		case protoreflect.DoubleKind:
			goType = "float64"
		case protoreflect.StringKind:
			goType = "string"
		case protoreflect.BytesKind:
			panic("bytes should be handled separately")
		case protoreflect.MessageKind, protoreflect.GroupKind:
			goType = field.Message.GoIdent.GoName
		case protoreflect.EnumKind:
			goType = field.Enum.GoIdent.GoName
		}
		return goType, pointer
	}

	var addField *Field
	var fieldPkgType string
	if field.Desc.IsList() {
		elementTyp, _ := getType(field)
		if field.Desc.Kind() == protoreflect.MessageKind {
			vt := gen.QualifiedGoIdent(field.Message.GoIdent)
			fieldPkgType = vt
			addField = m.AddMessageList(field.GoName, vt)
		} else {
			if field.Desc.Kind() == protoreflect.EnumKind {
				elementTyp = gen.QualifiedGoIdent(field.Enum.GoIdent)
				fieldPkgType = (elementTyp)
			}
			addField = m.AddList(field.GoName, elementTyp)
		}
	} else if field.Desc.IsMap() {
		keyField, valueField := field.Message.Fields[0], field.Message.Fields[1]
		keyTyp, _ := getType(keyField)
		valueTyp, _ := getType(valueField)
		if valueField.Desc.Kind() == protoreflect.MessageKind {
			vt := gen.QualifiedGoIdent(valueField.Message.GoIdent)
			fieldPkgType = (vt)
			addField = m.AddMessageMap(field.GoName, keyTyp, vt)
		} else {
			if valueField.Desc.Kind() == protoreflect.EnumKind {
				valueTyp = gen.QualifiedGoIdent(valueField.Enum.GoIdent)
				fieldPkgType = (valueTyp)
			}
			addField = m.AddMap(field.GoName, keyTyp, valueTyp)
		}
	} else if field.Desc.Kind() == protoreflect.BytesKind {
		addField = m.AddList(field.GoName, "byte")
	} else {
		// 基础类型 or 枚举 or 自定义类型
		elementTyp, opt := getType(field)
		if field.Desc.Kind() == protoreflect.MessageKind {
			vt := gen.QualifiedGoIdent(field.Message.GoIdent)
			fieldPkgType = (vt)
			addField = m.AddMessage(field.GoName, vt)
		} else {
			if field.Desc.Kind() == protoreflect.EnumKind {
				elementTyp = gen.QualifiedGoIdent(field.Enum.GoIdent)
				fieldPkgType = (elementTyp)
			}
			add := m.AddSimple
			if opt {
				add = m.AddPointer
			}
			addField = add(field.GoName, elementTyp)
		}
	}
	addField.Type = TypeString(addField, gen.ReadOnlyPkg)
	addField.Pkg, addField.OrgType = getTypePkg(fieldPkgType)
	addField.InnerName = "_" + addField.OuterName
}

// AddSimple 添加一个简单类型字段
func (m *Message) AddSimple(name, vt string) *Field {
	field := &Field{
		Type:      vt,
		OuterName: name,
		Kind:      KindSimple,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddPointer 添加一个指针字段
func (m *Message) AddPointer(name, vt string) *Field {
	field := &Field{
		Type:      vt,
		OuterName: name,
		Kind:      KindPointer,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddMessage 添加一个自定义类型字段
func (m *Message) AddMessage(name, vt string) *Field {
	field := &Field{
		Type:      pointerReadOnlyType(vt),
		OuterName: name,
		Kind:      KindMessage,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddList 添加一个列表字段
func (m *Message) AddList(name, vt string) *Field {
	field := &Field{
		Type:      vt,
		OuterName: name,
		Kind:      KindList,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddMessageList 添加一个自定义类型的列表字段
func (m *Message) AddMessageList(name, vt string) *Field {
	field := &Field{
		Type:      pointerReadOnlyType(vt),
		OuterName: name,
		Kind:      KindMessageList,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddMap 添加一个 map 字段
func (m *Message) AddMap(name, kt, vt string) *Field {
	field := &Field{
		Type:      mapType(kt, vt),
		OuterName: name,
		Kind:      KindMap,
	}
	m.Fields = append(m.Fields, field)
	return field
}

// AddMessageMap 添加一个自定义类型的 map 字段
func (m *Message) AddMessageMap(name, kt, vt string) *Field {
	field := &Field{
		Type:      mapType(kt, pointerReadOnlyType(vt)),
		OuterName: name,
		Kind:      KindMessageMap,
	}
	m.Fields = append(m.Fields, field)
	return field
}

func (m *Message) GenStruct(gen GenCfg) {
	gen.P("type ", m.Name, " struct {")
	gen.P("inner ", pointerType(m.PkgName))
	for _, field := range m.Fields {
		fs := field.FieldString(gen.ReadOnlyPkg)
		if fs == "" {
			continue
		}
		gen.P(fs)
	}
	gen.P("}")
}

func (m *Message) GenNew(gen GenCfg) {
	gen.P("func New", m.Name, "(p *", m.PkgName, ") ", pointerType(m.Name), " {")
	gen.P("if p == nil {")
	gen.P("return &", m.Name, "{inner: nil}")
	gen.P("}")
	gen.P("inner := p")
	gen.P("return &", m.Name, "{")
	gen.P("inner: inner,")
	for _, field := range m.Fields {
		if field.Kind == KindSimple {
			continue
		}
		gen.P(InitFieldString(field, gen.ReadOnlyPkg))
	}
	gen.P("}}")
	gen.P()
}

func (m *Message) GenGet(gen GenCfg) {
	for _, field := range m.Fields {
		gen.P("func (x *", m.Name, ") Get", field.OuterName, "() (_ ", field.Type, ") {")
		gen.P("if x == nil || x.inner == nil { return }")
		if field.Kind == KindSimple {
			gen.P(" return x.inner.Get", field.OuterName, "() }")
		} else {
			gen.P(" return x.", field.InnerName, "}")
		}
		gen.P()
	}
}

func genField(name, typ string) string {
	return strings.Join([]string{name, typ}, " ")
}

func mapType(kt, vt string) string {
	return strings.Join([]string{kt, vt}, ",")
}

func pointerReadOnlyType(typ string) string {
	return pointerType(readOnlyType(typ))
}

func pointerType(typ string) string {
	return "*" + typ
}

func readOnlyType(typ string) string {
	return typ + ReadOnlyMsgSuffix
}

func getTypePkg(ident string) (string, string) {
	if idx := strings.Index(ident, "."); idx != -1 {
		return ident[:idx], ident[idx+1:]
	}
	return "", ident
}
