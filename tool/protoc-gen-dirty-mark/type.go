package main

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"
)

type parsePlugin struct {
	Files     map[string]*File
	PBAlias   string
	MetaAlias string
	Imports   map[string]string // path -> alias
}

type File struct {
	Name      string
	Package   string
	PBAlias   string
	MetaAlias string
	Imports   map[string]string // path -> alias
	Structs   map[string]*Struct
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
	ret, err := format.Source(sb.Bytes())
	if err != nil {
		return nil, err
	}
	return ret, nil

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
func (s *Struct) AddValues(name, typ string) {
	s.Values = append(s.Values, Field{
		Type: typ,
		Name: name,
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
