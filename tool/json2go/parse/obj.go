package parse

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	FORMAT = "%s %s `json:\"%s,omitempty\"`"
)

var (
	translate = cases.Title(language.AmericanEnglish)
)

type JsonObject struct {
	KeyName  string
	TypeName string
	Fields   map[string]*JsonObject // 子字段
	Type     FiledType
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

func (j *JsonObject) GetCheckType() FiledType {
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
	}
	return typName
}

func (j *JsonObject) Key() string {
	if len(j.KeyName) == 0 {
		return "p"
	}
	return j.KeyName[:1]
}

func (j *JsonObject) FieldName() string {
	return title(j.KeyName)
}

func (j *JsonObject) ElemType() FiledType {
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
	var ret strings.Builder
	for _, s := range ss {
		ret.WriteString(translate.String(s))
	}
	return ret.String()
}
