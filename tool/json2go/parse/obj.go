package parse

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	FORMAT = "%s %s `json:\"%s,omitempty\"`"
)

var (
	translate = cases.Upper(language.AmericanEnglish)
)

type JsonObject struct {
	KeyName  string
	TypeName string
	MapType  string                 // 仅在对象类型是Map时有效.. 或者, 可以考虑不用这个?
	Fields   map[string]*JsonObject // 子字段
	Type     FieldType
}

func (j *JsonObject) AllFields() []*JsonObject {
	var ret = make([]*JsonObject, 0, len(j.Fields))
	for _, f := range j.Fields {
		ret = append(ret, f)
	}
	sort.Slice(ret, func(i, j2 int) bool {
		if ret[i].Type != ret[j2].Type {
			return ret[i].Type < ret[j2].Type
		}
		return ret[i].KeyName < ret[j2].KeyName
	})
	return ret
}

func (j *JsonObject) IsObject() bool {
	return j.Type.Check(TypeObject)
}

func (j *JsonObject) GetCheckType() FieldType {
	// Int不是标准的类型, 这里需要清空处理
	return j.Type.Clear(TypeInt | TypeFloat)
}

func (j *JsonObject) String() string {
	return fmt.Sprintf(FORMAT, j.FieldName(), "*"+j.ElemName(), j.KeyName)
}

func (j *JsonObject) Default() string {
	return j.Type.Default()
}

func (j *JsonObject) ElemName() string {
	var typName = j.Type.FiledType()
	if j.Type.Check(TypeObject) {
		typName += j.TypeName
	} else if j.Type.Check(TypeMap) {
		typName += j.MapType
	}
	return typName
}

func (j *JsonObject) TryCheckToMap() bool {
	if len(j.Fields) == 0 {
		return false
	}
	var checkKeyType = func(key string) FieldType {
		if _, err := strconv.ParseInt(key, 10, 64); err == nil {
			return TypeInt
		} else if _, err = strconv.ParseFloat(key, 64); err == nil {
			return TypeFloat
		}
		return 0
	}

	var keyType FieldType
	var valType FieldType
	var isMap = true
	for key, field := range j.Fields {
		if kt := checkKeyType(key); kt == 0 {
			return false
		} else {
			keyType |= kt
		}
		if valType == 0 {
			// 仅支持简单类型, 且必须做到类型一致值类型才可以
			valType = field.Type
			if !valType.Check(TypeBool | TypeFloat | TypeInt | TypeString) {
				isMap = false
				break
			}
		} else {
			if valType.NaiveType() != field.Type.NaiveType() {
				isMap = false
				break
			}
		}
	}
	if !isMap {
		return false
	}
	// 包装成Map类型
	var key = keyType.FiledType()
	var val = valType.FiledType()
	j.MapType = fmt.Sprintf("map[%s]%s", key, val)
	j.Type.SetMap()
	return true
}

func (j *JsonObject) Key() string {
	return "p"
}

func (j *JsonObject) FieldName() string {
	return title(j.KeyName)
}

func (j *JsonObject) ElemType() FieldType {
	return j.Type.ElemType()
}

func (j *JsonObject) Merge(src *JsonObject) error {
	for name, obj := range src.Fields {
		if old, ok := j.Fields[name]; ok {
			if old.GetCheckType() != obj.GetCheckType() {
				return fmt.Errorf("check %v type not match", name)
			}
			if old.GetCheckType() == TypeObject {
				if err := old.Merge(obj); err != nil {
					return fmt.Errorf("merge %v error:%w", name, err)
				}
			} else {
				old.Type |= obj.Type.NumberType()
			}
		} else {
			j.Fields[name] = obj
		}
	}
	return nil
}

func title(src string) string {
	var ss = strings.Split(src, "_")
	var ret = bytes.NewBuffer(nil)
	for _, s := range ss {
		ret.WriteString(strings.Title(s))
	}
	//TODO 处理非法的标识符
	var bs = ret.Bytes()
	if len(bs) == 0 {
		return ""
	}
	var first = bs[0]
	if first >= 'A' && first <= 'Z' {
		return ret.String()
	}
	return "F" + ret.String()
}
