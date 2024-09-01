package main

const FileTemplate = `package {{.Package}}

import (
	{{- range $path, $alias := .Imports}}
		{{- if $alias}}
	{{$alias}} "{{$path}}"
		{{- else}}
	"{{$path}}"
		{{- end}}
	{{- end}}
)


{{- $top := .}}

{{- range $name, $struct := .Structs}}

const (
	_{{$name}}FieldIndex = iota - 1
	{{- range $field := $struct.AllFields}}
	{{$name}}FieldIndex{{$field.Name}}
	{{- end}}
	{{$name}}FieldMax
)

var _{{$name}}ApplyDirtyTable = []func(*{{$name}}, *{{$top.PBAlias}}.{{$name}}){
	{{- range $field := $struct.AllFields}}
	{{$name}}FieldIndex{{$field.Name}}: (*{{$name}}).applyDirty{{$field.Name}},
	{{- end}}
}

type {{$name}} struct {
	mark *{{$top.MetaAlias}}.BitsetMark
	{{- range $field := $struct.Values}}
	_{{$field.Name}} {{$field.Type}}
	{{- end}}
	{{- range $field := $struct.References}}
	_{{$field.Name}} {{$field.Type}}
	{{- end}}
	{{- range $field := $struct.Containers}}
	_{{$field.Name}} *{{$top.MetaAlias}}.{{$field.Extra}}{{$field.Type}}
	{{- end}}
}

func New{{$name}}() *{{$name}} {
	var m {{$name}}
	m.mark = {{$top.MetaAlias}}.NewBitsetMark({{$name}}FieldMax)
	{{- range $field := $struct.References}}
	m._{{$field.Name}} = New{{$field.Name}}()
	{{- end}}
	{{- range $field := $struct.Containers}}
	m._{{$field.Name}} = {{$top.MetaAlias}}.New{{$field.Extra}}{{$field.Type}}()
	{{- end}}
	return &m
}

// NewValue creates a new {{$name}}.
func (*{{$name}}) NewValue() {{$top.MetaAlias}}.IValue[*{{$top.PBAlias}}.{{$name}}] {
	return New{{$name}}()
}


{{- range $field := $struct.Values}}
// Get{{$field.Name}} gets the {{$field.Name}}.
func (m *{{$name}}) Get{{$field.Name}}() {{$field.Type}} {
	return m._{{$field.Name}}
}

// Set{{$field.Name}} sets the {{$field.Name}}.
func (m *{{$name}}) Set{{$field.Name}}(v {{$field.Type}}) {
	m._{{$field.Name}} = v
	{{$top.MetaAlias}}.MarkHelper(m.mark, {{$name}}FieldIndex{{$field.Name}})
}


func (m *{{$name}}) applyDirty{{$field.Name}}(p *{{$top.PBAlias}}.{{$name}}) {
	{{- if eq $field.Extra "optional" -}}
	p.{{$field.Name}} = {{$top.MetaAlias}}.Pointer(m.Get{{$field.Name}}())
	{{- else -}}
	p.{{$field.Name}} = m.Get{{$field.Name}}()
	{{- end}}
}
{{- end}}
{{- range $field := $struct.References}}
// Get{{$field.Name}} gets the {{$field.Name}}.
func (m *{{$name}}) Get{{$field.Name}}() {{$field.Type}} {
	if m._{{$field.Name}} == nil {
		m._{{$field.Name}} = New{{$field.Name}}()
	}
	{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}.mark, m.mark, {{$name}}FieldIndex{{$field.Name}}) 
	return m._{{$field.Name}}
}

// Set{{$field.Name}} sets the {{$field.Name}}.
func (m *{{$name}}) Set{{$field.Name}}(v {{$field.Type}}) {
	if m._{{$field.Name}} == v {
		return
	}
	if m._{{$field.Name}} != nil {
		{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}.mark, nil, 0) 
	}
	m._{{$field.Name}} = v
	if v != nil {
		{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}.mark, m.mark, {{$name}}FieldIndex{{$field.Name}}) 
	}
	{{$top.MetaAlias}}.MarkHelper(m.mark, {{$name}}FieldIndex{{$field.Name}})
}

func (m *{{$name}}) applyDirty{{$field.Name}}(p *{{$top.PBAlias}}.{{$name}}) {
	if p.{{$field.Name}} == nil {
		p.{{$field.Name}} = &{{$top.PBAlias}}.{{$field.Name}}{}
	}
	m.Get{{$field.Name}}().DirtyCollect(p.{{$field.Name}})
}
{{- end}}
{{- range $field := $struct.Containers}}
// Get{{$field.Name}} gets the {{$field.Name}}.
func (m *{{$name}}) Get{{$field.Name}}() *{{$top.MetaAlias}}.{{$field.Extra}}{{$field.Type}} {
	if m._{{$field.Name}} == nil {
		m._{{$field.Name}} = {{$top.MetaAlias}}.New{{$field.Extra}}{{$field.Type}}()
	}
	{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}, m.mark, {{$name}}FieldIndex{{$field.Name}}) 
	return m._{{$field.Name}}
}

// Set{{$field.Name}} sets the {{$field.Name}}.
func (m *{{$name}}) Set{{$field.Name}}(v *{{$top.MetaAlias}}.{{$field.Extra}}{{$field.Type}}) {
	if m._{{$field.Name}} == v {
		return
	}
	if m._{{$field.Name}} != nil {
		{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}, nil, 0) 
	}
	m._{{$field.Name}} = v
	if v != nil {
		{{$top.MetaAlias}}.SetMarkHelper(m._{{$field.Name}}, m.mark, {{$name}}FieldIndex{{$field.Name}})
	}
	{{$top.MetaAlias}}.MarkHelper(m.mark, {{$name}}FieldIndex{{$field.Name}})
}


func (m *{{$name}}) applyDirty{{$field.Name}}(p *{{$top.PBAlias}}.{{$name}}) {
	p.{{$field.Name}} = m.Get{{$field.Name}}().DirtyCollect(p.{{$field.Name}})
}
{{- end}}


// FromProto sets the value from the target.
func (m *{{$name}}) FromProto(p *{{$top.PBAlias}}.{{$name}}) {
	{{- range $field := $struct.Values}}
	m.Set{{$field.Name}}(p.Get{{$field.Name}}())
	{{- end}}
	{{- range $field := $struct.References}}
	m.Get{{$field.Name}}().FromProto(p.Get{{$field.Name}}())
	{{- end}}
	{{- range $field := $struct.Containers}}
	m.Get{{$field.Name}}().FromProto(p.Get{{$field.Name}}())
	{{- end}}
}

// ToProto gets the target from the value.
func (m *{{$name}}) ToProto() *{{$top.PBAlias}}.{{$name}} {
	var p {{$top.PBAlias}}.{{$name}}
	{{- range $field := $struct.Values}}
		{{- if eq $field.Extra "optional"}}
	p.{{$field.Name}} = {{$top.MetaAlias}}.Pointer(m.Get{{$field.Name}}())
		{{- else }}
	p.{{$field.Name}} = m.Get{{$field.Name}}()
		{{- end}}
	{{- end}}
	{{- range $field := $struct.References}}
	p.{{$field.Name}} = m.Get{{$field.Name}}().ToProto()
	{{- end}}
	{{- range $field := $struct.Containers}}
	p.{{$field.Name}} = m.Get{{$field.Name}}().ToProto()
	{{- end}}
	return &p
}

// resetDirty resets the dirty mark.
func (m *{{$name}}) resetDirty() {
	{{$top.MetaAlias}}.ResetHelper(m.mark)
	{{- range $field := $struct.References}}
	m.Get{{$field.Name}}().resetDirty()
	{{- end}}
}

// DirtyProto returns proto apply the dirty mark.
func (m *{{$name}}) DirtyProto() *{{$top.PBAlias}}.{{$name}} {
	var p {{$top.PBAlias}}.{{$name}}
	m.DirtyCollect(&p)
	return &p
}


// DirtyCollect applies the dirty mark to the target.
func (m *{{$name}}) DirtyCollect(target *{{$top.PBAlias}}.{{$name}}) {
	for dirtyIdx := range {{$top.MetaAlias}}.DirtyBitsHelper(m.mark) {
		_{{$name}}ApplyDirtyTable[dirtyIdx](m, target)
	}
	m.resetDirty()
}

// GetMark gets the mark.
func (m *{{$name}}) GetMark() {{$top.MetaAlias}}.IMark {
	return m.mark
}

{{- end}}

`
