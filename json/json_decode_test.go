package json

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyNumber string

func (m *MyNumber) UnmarshalJSON(bytes []byte) error {
	log.Printf("data:%s", bytes)

	if len(bytes) == 0 {
		*m = ""
		return nil
	}
	if len(bytes) >= 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		*m = MyNumber(string(bytes[1 : len(bytes)-1]))
	} else {
		*m = MyNumber(string(bytes))
	}
	return nil
}

func TestJsonDecode(t *testing.T) {
	var data1 struct {
		Number int `json:"num"`
	}

	var data2 struct {
		Number string `json:"num"`
	}

	type data3 struct {
		Number MyNumber `json:"num"`
	}

	var testCases = []struct {
		V1 int
		V2 string
	}{
		{
			V1: 0,
			V2: "0",
		},
		{
			V1: 1,
			V2: "1",
		},
		{
			V1: 100,
			V2: "100",
		},
		{
			V1: 0,
			V2: "干干干",
		},
	}

	for _, tc := range testCases {
		data1.Number = tc.V1
		data2.Number = tc.V2
		var bs1, _ = json.Marshal(data1)
		var bs2, _ = json.Marshal(data2)
		var d1, d2 data3
		_ = json.Unmarshal(bs1, &d1)
		_ = json.Unmarshal(bs2, &d2)
		assert.Equal(t, d1.Number, d2.Number)
	}

	var bs, _ = json.Marshal(data3{})
	t.Logf("%s", bs)
}

func TestJsonStruct(t *testing.T) {
	// 定义json结构体时, 0值尽量不要做为一个有效值, 空值可能和零值发生混淆
	// 如果非要做为一个有效值, 可以考虑将其设置为对应类型的指针+omitempty

	// var val string

	type Data struct {
		D *string `json:"data,omitempty"`
	}

	var src1 = []byte(`{}`)
	var src2 = []byte(`{"data":"1231"}`)

	var d1, d2 Data
	_ = json.Unmarshal(src1, &d1)
	_ = json.Unmarshal(src2, &d2)

	t.Logf("%+v, %+v", d1.D, *d2.D)
}
