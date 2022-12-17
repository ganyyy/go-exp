package main

import (
	"reflect"
	"testing"

	"ganyyy.com/go-exp/rpc/grpc/proto"
	"github.com/stretchr/testify/assert"
)

type dataType struct {
	_   uintptr
	typ uintptr
}

func getType(v interface{}) reflect.Type {
	var rtype = reflect.TypeOf(v)
	return rtype
}

func TestShowType(t *testing.T) {
	var p1 *proto.RepeatData_Data_Req
	var p2 *proto.RepeatData_Data_Rsp

	t.Logf("p1:%v", getType(p1))
	t.Logf("p2:%v", getType(p2))

	var a = &proto.RepeatData_Data{
		Data: &proto.RepeatData_Data_Req{},
	}
	t.Logf("a:%v", getType(a.GetData()))

	assert.Equal(t, getType(p1), getType(a.GetData()))
	var b = &proto.RepeatData_Data{
		Data: &proto.RepeatData_Data_Rsp{},
	}
	t.Logf("b:%v", getType(b.GetData()))
	assert.Equal(t, getType(p2), getType(b.GetData()))
}

type genHandle func(*proto.RepeatData_Data) Handle

var (
	registerMap = make(map[reflect.Type]genHandle)
)

func register(t interface{}, handle genHandle) {
	registerMap[reflect.TypeOf(t)] = handle
}

func TestHandleType(t *testing.T) {
	register((*proto.RepeatData_Data_Req)(nil), func(rd *proto.RepeatData_Data) Handle { return &HandleReq{HelloRequest: rd.GetReq()} })
	register((*proto.RepeatData_Data_Rsp)(nil), func(rd *proto.RepeatData_Data) Handle { return &HandleRsp{HelloResponse: rd.GetRsp()} })

	show := func(data *proto.RepeatData_Data) {
		var gen = registerMap[reflect.TypeOf(data.GetData())]
		assert.NotNil(t, gen)
		req := gen(data)
		t.Logf("%#v, %v", req, req)
	}

	show(&proto.RepeatData_Data{Data: &proto.RepeatData_Data_Req{Req: &proto.HelloRequest{Name: "123"}}})
	show(&proto.RepeatData_Data{Data: &proto.RepeatData_Data_Rsp{Rsp: &proto.HelloResponse{Message: "hello"}}})
}
