package main

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ImportInfo struct {
	PBAlias   string
	MetaAlias string
	Imports   map[string]string // path -> alias
}

type parsePlugin struct {
	Files map[string]*File
	ImportInfo
}

type File struct {
	*ImportInfo
	Name    string
	Package string
	Structs map[string]*Struct // name -> struct
}

// Render renders the file.
func (f *File) Render() ([]byte, error) {
	var sb bytes.Buffer
	tp, err := template.New(f.Name).Parse(FileTemplate)
	if err != nil {
		return nil, err
	}
	err = tp.Execute(&sb, f)
	if err != nil {
		return nil, err
	}
	// return sb.Bytes(), nil
	ret, err := format.Source(sb.Bytes())
	if err != nil {
		return nil, err
	}
	return ret, nil

}

// Write writes the file to the given path.
func (f *File) Write(path string) error {
	content, err := f.Render()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, f.Name), content, 0644)
}

type Field struct {
	Type  string
	Name  string
	Extra string
}

type Struct struct {
	Name       string
	Values     []Field // basic types
	References []Field // reference types
	Containers []Field // map or list
}

// AddValues adds a value to the struct.
func (s *Struct) AddValues(name, typ string, extra ...string) {
	s.Values = append(s.Values, Field{
		Type:  typ,
		Name:  name,
		Extra: strings.Join(extra, " "),
	})
}

func ReferenceType(name string) string {
	return "*" + name
}

func JoinGenerics(generics ...string) string {
	return "[" + strings.Join(generics, ", ") + "]"
}

func PackageType(pkg, name string) string {
	return pkg + "." + name
}

// AddReferences adds a reference to the struct.
func (s *Struct) AddReferences(name string) {
	s.References = append(s.References, Field{
		Name: name,
		Type: ReferenceType(name),
	})
}

// AddValuesMap adds a value map to the struct.
func (s *Struct) AddValuesMap(name, kt, vt string) {
	s.Containers = append(s.Containers, Field{
		Type:  JoinGenerics(kt, vt),
		Name:  name,
		Extra: "ValueMap",
	})
}

// AddReferencesMap adds a reference map to the struct.
func (s *Struct) AddReferencesMap(name, kt, vt string) {
	s.Containers = append(s.Containers, Field{
		Type:  JoinGenerics(kt, ReferenceType(vt), ReferenceType(PackageType(*pbAlias, vt))),
		Name:  name,
		Extra: "ReferenceMap",
	})
}

// AddValuesList adds a value list to the struct.
func (s *Struct) AddValuesList(name, vt string) {
	s.Containers = append(s.Containers, Field{
		Type:  JoinGenerics(vt),
		Name:  name,
		Extra: "ValueList",
	})
}

// AddReferencesList adds a reference list to the struct.
func (s *Struct) AddReferencesList(name, vt string) {
	s.Containers = append(s.Containers, Field{
		Type:  JoinGenerics(ReferenceType(vt), ReferenceType(PackageType(*pbAlias, vt))),
		Name:  name,
		Extra: "ReferenceList",
	})
}

// AllFields returns all fields.
func (s *Struct) AllFields() []Field {
	var fields = make([]Field, 0, len(s.Values)+len(s.References)+len(s.Containers))
	fields = append(fields, s.Values...)
	fields = append(fields, s.References...)
	fields = append(fields, s.Containers...)
	return fields
}

func parseField(field *protogen.Field, s *Struct) {
	var getType = func(field *protogen.Field) (string, bool) {
		var goType string
		var pointer = field.Desc.HasPresence() // optional fields are pointers
		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			goType = "bool"
		case protoreflect.EnumKind:
			// TODO handle enum
			goType = field.Enum.GoIdent.GoName
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
		}
		return goType, pointer
	}

	if field.Desc.IsMap() {
		// current not support map[x][]byte
		keyField, valueField := field.Message.Fields[0], field.Message.Fields[1]
		keyName, _ := getType(keyField)
		valName, _ := getType(valueField)
		if valueField.Desc.Kind() == protoreflect.MessageKind {
			s.AddReferencesMap(field.GoName, keyName, valName)
		} else {
			s.AddValuesMap(field.GoName, keyName, valName)
		}
	} else if field.Desc.IsList() {
		// current not support [][]byte
		element, _ := getType(field)
		if field.Desc.Kind() == protoreflect.MessageKind {
			s.AddReferencesList(field.GoName, element)
		} else {
			s.AddValuesList(field.GoName, element)
		}
	} else if field.Desc.Kind() == protoreflect.BytesKind {
		s.AddValuesList(field.GoName, "byte")
	} else {
		element, pointer := getType(field)
		if field.Desc.Kind() == protoreflect.MessageKind {
			s.AddReferences(field.GoName)
		} else {
			var extra string
			if pointer {
				extra = "optional"
			}
			s.AddValues(field.GoName, element, extra)
		}
	}
}
