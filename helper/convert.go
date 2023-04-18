package helper

import (
	"encoding/json"
	"os"

	"google.golang.org/protobuf/proto"
)

func To[T any](v any) T {
	var ret, ok = v.(T)
	if !ok {
		var rr T
		return rr
	}
	return ret
}

type Data interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type JsonStruct struct{}

func (d JsonStruct) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (d JsonStruct) Unmarshal(bs []byte, v interface{}) error {
	return json.Unmarshal(bs, v)
}

type ProtoStruct struct{}

func (p ProtoStruct) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p ProtoStruct) Unmarshal(bs []byte, v interface{}) error {
	return proto.Unmarshal(bs, v.(proto.Message))
}

func Load[T any](data Data, path string) (v T, err error) {
	bs, _ := os.ReadFile(path)
	err = data.Unmarshal(bs, &v)
	return
}

var (
	jsonData  = JsonStruct{}
	protoData = ProtoStruct{}
)

var _ = func() {
	type S struct {
		Name string
	}
	var t, e = Load[S](jsonData, "path")
	_ = e
	_ = t.Name
	t, e = Load[S](protoData, "path")
	_ = e
	_ = t.Name
}
