package logic

import (
	"protoc-gen-go-handle/gen/proto"
	"testing"

	gp "google.golang.org/protobuf/proto"
)

func TestHandle(t *testing.T) {
	var req proto.GetPlayerInfoReq
	req.Acid = "123"
	data, _ := gp.Marshal(&req)
	resp := Handle(int32(proto.LogicCode_GetPlayerInfo), data)
	var rsp proto.GetPlayerInfoRsp
	_ = gp.Unmarshal(resp, &rsp)
	t.Logf("Req:%v, Rsp:%v", &req, &rsp)
}
