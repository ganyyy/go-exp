package nc

import "encoding/json"

//TODO 尝试自己搞一个编解码器

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type jsonCodec struct{}

func (_ jsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (_ jsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var JsonCodec jsonCodec
