package logic

import (
	"protoc-gen-go-handle/gen/proto"
	"reflect"

	gp "google.golang.org/protobuf/proto"
)

type Player struct {
}

type IRequest interface {
	GetReq() gp.Message
	GetRsp() gp.Message
	Handle(*Player)
}

var (
	requestMap = make(map[proto.LogicCode]reflect.Type)
)

func init() {
	requestMap[proto.LogicCode_GetPlayerInfo] = reflect.TypeOf(GetPlayerInfo{})
}

func Handle(code int32, data []byte) []byte {
	requestType, ok := requestMap[proto.LogicCode(code)]
	if !ok {
		return nil
	}
	handle, ok := reflect.New(requestType).Interface().(IRequest)
	if !ok {
		return nil
	}
	gp.Unmarshal(data, handle.GetReq())
	handle.Handle(&Player{})
	rsp, _ := gp.Marshal(handle.GetRsp())
	return rsp
}
