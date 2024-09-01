package main

import (
	"os"
	"protoc-gen-dirty-mark/data"
	"protoc-gen-dirty-mark/meta"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	var file File
	file.ImportInfo = &ImportInfo{}

	var inner Struct
	inner.Name = "Inner"
	inner.AddValues("Data", "string")
	inner.AddValues("Age", "int32")

	*pbAlias = "pb123"
	*metaAlias = "meta1"

	var data Struct
	data.Name = "Data"
	data.AddValues("Name", "string", "optional")
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

	content, err := file.Render()

	require.NoError(t, err)
	_, err = f.Write(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

func TestData(t *testing.T) {
	var dataObj = data.NewData()
	dataObj.SetName("test")
	innerObj := dataObj.GetInner()
	innerObj.SetData("inner")

	logDirty := func() {

		var pc [1]uintptr
		n := runtime.Callers(2, pc[:])
		if n == 1 {
			frames := runtime.CallersFrames(pc[:])
			frame, _ := frames.Next()
			t.Logf("Called from %s:%d", frame.Function, frame.Line)
		}

		dirty := dataObj.DirtyProto()
		t.Logf("%+v", dirty.String())
	}

	strList := dataObj.GetStrList()

	strList.Add("str1")
	strList.Add("str2")
	logDirty()
	strs := strList.ToProto()
	t.Logf("%+v", strs)

	dataObj.SetName("test2")
	logDirty()

	strMap := meta.NewValueMap[string, string]()
	strMap.Set("key", "value")
	strMap.Set("key2", "value2")
	dataObj.SetStrMap(strMap)
	logDirty()

	innerList := meta.NewReferenceList[*data.Inner]()
	dataObj.SetInnerList(innerList)
	dataObj.SetInner(data.NewInner())
	innerList.Add(innerObj)
	logDirty()
	innerObj, ok := innerList.Remove(0)
	require.True(t, ok)
	innerList.Insert(0, innerObj)
	logDirty()

	dataObj.GetInnerMap()

	strMap.Del("key")
	logDirty()

	innerObj.SetData("inner2")
	logDirty()
}
