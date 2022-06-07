package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type (
	naiveObj   map[string]any
	naiveValue any
	naiveSlice []any
)

var (
	ErrEmptyObj    = errors.New("empty object")
	ErrInvalidType = errors.New("invalid type")
)

func parseObject(src naiveObj) (*JsonObject, error) {
	if src == nil {
		return nil, ErrEmptyObj
	}
	var obj = new(JsonObject)
	obj.Type = TypeObject
	obj.Fields = make(map[string]*JsonObject)
	for name, val := range src {
		var field, err = parseValue(val)
		if err != nil {
			return nil, fmt.Errorf("parse %v, %+v error:%w", name, val, err)
		}
		field.KeyName = name
		obj.Fields[name] = field
	}
	return obj, nil
}

//parseValue current: 用来转义数组时使用的同级别结构, 单个结构体可以直接跳过
func parseValue(src naiveValue) (*JsonObject, error) {
	if src == nil {
		return nil, ErrEmptyObj
	}
	var vt = reflect.TypeOf(src)
	switch vt.Kind() {
	case reflect.Slice:
		return parseSlice(src.([]any))
	case reflect.Bool:
		return &JsonObject{Type: TypeBool}, nil
	case reflect.String:
		// 如果上层使用了UseNumber, 所以这里需要判断一下是数字类型还是普通的字符串
		if num, ok := src.(json.Number); ok {
			var rt = TypeFloat
			if _, err := num.Int64(); err == nil {
				rt = TypeInt
			}
			return &JsonObject{Type: rt}, nil
		} else {
			return &JsonObject{Type: TypeString}, nil
		}
	case reflect.Float64:
		// 如果未开启UseNumber的情况下, 只有float类型
		return &JsonObject{Type: TypeFloat}, nil
	case reflect.Map:
		return parseObject(src.(map[string]any))
	case reflect.Invalid:
		// 这是啥? null 吗?
		fallthrough
	default:
		// 不认识的类型?
		return nil, ErrInvalidType
	}

}

func parseSlice(src naiveSlice) (*JsonObject, error) {
	var current *JsonObject
	var err error
	for _, obj := range src {
		var parseVal *JsonObject
		parseVal, err = parseValue(obj)
		if err != nil {
			return nil, fmt.Errorf("parse %+v error:%w", obj, err)
		}
		if current == nil {
			current = parseVal
		} else {
			if current.GetCheckType() != parseVal.GetCheckType() {
				return nil, fmt.Errorf("parse slice found not match type %+v", obj)
			}
			// Object 校验字段
			if current.GetCheckType().Check(TypeObject) {
				if err := current.Merge(parseVal); err != nil {
					return nil, err
				}
			} else {
				// 特殊处理以下数字类型
				current.Type |= parseVal.Type.NumberType()
			}
		}
	}
	if current.Type.Check(TypeSlice) {
		current.Type.AddSlice(1)
	} else {
		current.Type.SetSlice()
	}
	return current, nil
}

func ParseAllType(root *JsonObject) []*JsonObject {
	var allType []*JsonObject
	//把自己先加进去
	allType = append(allType, root)
	var dfs func(obj *JsonObject)
	dfs = func(jt *JsonObject) {
		for name, obj := range jt.Fields {
			if obj.Type&TypeObject != 0 {
				obj.TypeName = jt.TypeName + title(name)
				dfs(obj)
				allType = append(allType, obj)
			}
		}
	}
	dfs(root)

	sort.Slice(allType, func(i, j int) bool {
		return allType[i].TypeName < allType[j].TypeName
	})

	return allType
}
