package parse

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testParseJson(t *testing.T, src string, fileName string, base string) {
	var decoder = json.NewDecoder(strings.NewReader(src))
	// decoder.UseNumber()
	var v naiveValue
	assert.Nil(t, decoder.Decode(&v))

	var obj, err = parseValue(v)
	if !assert.Nil(t, err) {
		return
	}
	obj.KeyName = base
	obj.TypeName = title(obj.KeyName)

	if obj.Type.Check(TypeSlice) {
		t.Logf("slice count:%v", obj.Type.SliceCount())
		t.Logf("slice type:%v", obj.ElemName())
	}

	var allType = ParseAllType(obj)
	for _, v := range allType {
		t.Logf("type:%v", v.TypeName)
		for _, f := range v.AllFields() {
			t.Logf("\t%+v", f)
		}
	}
	assert.Nil(t, err)

	var output = "./data/" + fileName
	var tp TemplateParse
	tp.Root = obj
	tp.AllType = allType
	tp.PkgName = "data"
	err = tp.Parse(output)
	t.Log(err)
}

const (
	parse_data1 = `
	[
    {
        "name": "123",
        "age": 10.2,
        "list": [
            {
                "name": "123",
                "age": 100,
                "qq": "7788"
            },
            {
                "name": "456",
                "age": 200.2,
                "city": "7788"
            }
        ]
    },
    {
        "name": "123",
        "age": 10,
        "list2": [
            {
                "name": "123",
                "age": 100,
                "qq": "7788"
            },
            {
                "name": "456",
                "age": 200.2,
                "city": "7788"
            }
        ],
        "other": {
            "is_true": true,
            "add": [
                100,
                200,
                300,
                400
            ]
        }
    }
]
	`

	parse_data2 = `
	[
		[
			[
				{
					"name":"123",
					"age":123
				}
			],
			[
				{
					"addr":"123",
					"other":123
				}
			]
		]	
	]
	`
	parse_data3 = `
	{
		"name":"123",
		"age":123
	}
	`

	parse_data4 = `
	[
		1,2,3,4
	]
	`
)

func TestParse(t *testing.T) {
	var param ParseParam
	param.GoPackage = "data"
	param.OutputPath = "./data"
	param.UseNumber = true

	if err := param.InitOutput(); err != nil {
		t.Logf("init error:%v", err)
		t.FailNow()
	}

	t.Run("parse", func(t *testing.T) {
		t.Log(parseInputData("data", []byte(parse_data1), &param))
	})

	t.Run("parse2", func(t *testing.T) {
		t.Log(parseInputData("user", []byte(parse_data2), &param))
	})

	t.Run("parse3", func(t *testing.T) {
		t.Log(parseInputData("object", []byte(parse_data3), &param))
	})

	t.Run("parse4", func(t *testing.T) {
		t.Log(parseInputData("float_data", []byte(parse_data4), &param))
	})
}

func TestFileTypeSlice(t *testing.T) {
	var ft FiledType
	ft.SetSlice()
	ft.AddSlice(10)
	t.Logf("%032b", ft)
	t.Logf("%v", ft.FiledType())
}
