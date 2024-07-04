package main

import (
	"os"
	"protoc-gen-dirty-mark/data"
	"protoc-gen-dirty-mark/meta"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	var file File

	var inner Struct
	inner.Name = "Inner"
	inner.AddValues("Data", "string")

	*pbAlias = "pb123"
	*metaAlias = "meta1"

	var data Struct
	data.Name = "Data"
	data.AddValues("Name", "string")
	data.AddReferences("Inner")
	data.AddValuesMap("StrMap", "string", "string")
	data.AddReferencesMap("InnerMap", "string", "Inner")
	data.AddValuesList("StrList", "string")
	data.AddReferencesList("InnerList", "Inner")

	file.Package = "data"
	file.MetaAlias = *metaAlias
	file.PBAlias = *pbAlias
	file.Imports = map[string]string{
		"protoc-gen-dirty-mark/pb":   file.PBAlias,
		"protoc-gen-dirty-mark/meta": file.MetaAlias,
	}
	file.Structs = map[string]*Struct{
		"Inner": &inner,
		"Data":  &data,
	}

	var f, err = os.OpenFile("data/data.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(t, err)

	tp, err := template.New("data").Parse(FileTemplate)
	require.NoError(t, err)

	err = tp.Execute(f, &file)
	require.NoError(t, err)
}

func TestData(t *testing.T) {
	var d = data.NewData()
	d.SetName("test")
	i := data.NewInner()
	i.SetData("inner")
	d.SetInner(i)

	logDirty := func() {
		dirty := d.DirtyProto()
		t.Logf("%+v", dirty.String())
	}

	logDirty()
	d.SetName("test2")
	logDirty()

	strMap := meta.NewValueMap[string, string]()
	strMap.Set("key", "value")
	strMap.Set("key2", "value2")
	d.SetStrMap(strMap)
	logDirty()

	strMap.Del("key")
	logDirty()
}
