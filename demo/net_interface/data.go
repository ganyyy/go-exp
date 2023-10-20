package netinterface

import (
	"encoding/json"
	"errors"
)

var (
	ErrUnknownDataType = errors.New("unknown data type")
)

type NetData struct {
	Type DataType
	Data json.RawMessage
}

func MarshalDataToNet(data IData) ([]byte, error) {
	meta, ok := dataTypeToMeta[data.GetType()]
	if !ok {
		return nil, ErrUnknownDataType
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&NetData{
		Type: meta.GetType(),
		Data: bs,
	})
}

func UnmarshalDataFromNet[T IData](bs []byte) (T, error) {
	var netData NetData
	var empty T
	if err := json.Unmarshal(bs, &netData); err != nil {
		return empty, err
	}
	meta, ok := dataTypeToMeta[netData.Type]
	if !ok {
		return empty, ErrUnknownDataType
	}
	data := meta.New()
	if err := json.Unmarshal(netData.Data, data); err != nil {
		return empty, err
	}
	return data.(T), nil
}
