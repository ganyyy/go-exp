package parse

import (
	"testing"
)

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
		1,2,3,4.5
	]
	`

	parse_data5 = `
	[
		[
			1,2,3
		],
		[

		]
	]
	`

	parse_data6 = `
	{
		"configMap":{
			"100": 10,
			"101": 10,
			"102": 10
		}
	}
	`

	parse_data7 = `
	[
		[
			{
				"100": 10,
				"101": 10,
				"102": 10
			}
		],
		[
			{
				"100": 10,
				"101": 10,
				"102": 10
			}
		],
		[
			{
				"100": 10,
				"101": 10,
				"102": 10
			}
		]
	]
	`
)

func TestParse(t *testing.T) {
	var param ParseParam
	param.GoPackage = "data"
	param.OutputPath = "./data"
	param.UseNumber = true
	param.ParseMap = true

	if err := param.InitOutput(); err != nil {
		t.Logf("init error:%v", err)
		t.FailNow()
	}

	var cases = []struct {
		name string
		data string
	}{
		{"data", parse_data1},
		{"user", parse_data2},
		{"object", parse_data3},
		{"float_data", parse_data4},
		{"empty_slice5", parse_data5},
		{"single_map", parse_data6},
		{"multi_map", parse_data7},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Log(parseInputData(c.name, []byte(c.data), &param))
		})
	}
}

func TestFileTypeSlice(t *testing.T) {
	var ft FieldType
	ft.SetSlice()
	ft.AddSlice(10)
	t.Logf("%032b", ft)
	t.Logf("%v", ft.FiledType())
}
