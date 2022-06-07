package parse

import (
	"bytes"
	"go/format"
	"html/template"
	"io/fs"
	"os"
)

const templateFormat = `
package {{.PkgName}}

import (
	"encoding/json"
	"os"
)


var src{{.Root.TypeName}} {{.Root.ElemName}}

func Load{{.Root.TypeName}}(path string) error {
	var bs, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	return Parse{{.Root.TypeName}}(bs)
}

func Parse{{.Root.TypeName}}(bs []byte) error {
	return json.Unmarshal(bs, &src{{.Root.TypeName}})
}

func Get{{.Root.TypeName}} () {{.Root.ElemName}} {
	return src{{.Root.TypeName}}
}

{{range $type := .AllType}}
{{if $type.IsObject}}
{{$typeName := $type.TypeName}}
type {{$typeName}} struct {
    {{range $field := $type.AllFields}}
    {{$field.String | raw}}{{end}}
}

{{range $field := $type.AllFields}}
    {{$key := $type.Key}}
    {{$keyName := $field.FieldName}}
func({{$key}} *{{$typeName}}) Get{{$keyName}}() {{$field.ElemName}} {
    if {{$key}} == nil || {{$key}}.{{$keyName}} == nil {
        return {{$field.Default | raw}}
    }
    return *{{$key}}.{{$keyName}}
}
{{end}}
{{end}}
{{end}}
`

var t *template.Template

func init() {
	var err error
	t, err = template.New("root").Funcs(template.FuncMap{
		"raw": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).Parse(templateFormat)
	if err != nil {
		panic(err)
	}
}

type TemplateParse struct {
	PkgName string
	Root    *JsonObject
	AllType []*JsonObject
}

func (p *TemplateParse) Parse(output string) error {
	var outBuf = bytes.NewBuffer(nil)
	// 翻译
	if err := t.Execute(outBuf, p); err != nil {
		return err
	}

	// go fmt 格式化
	content, err := format.Source(outBuf.Bytes())
	if err != nil {
		os.WriteFile(output, outBuf.Bytes(), fs.ModePerm)
		return err
	}
	// 写入到目标文件
	return os.WriteFile(output, content, fs.ModePerm)
}
