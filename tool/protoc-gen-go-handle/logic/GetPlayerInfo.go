package logic

import (
	"protoc-gen-go-handle/gen/proto"

	gp "google.golang.org/protobuf/proto"
)

type GetPlayerInfo struct {
	Req proto.GetPlayerInfoReq
	Rsp proto.GetPlayerInfoRsp
}

func (p *GetPlayerInfo) GetReq() gp.Message {
	return &p.Req
}

func (p *GetPlayerInfo) GetRsp() gp.Message {
	return &p.Rsp
}

func (p *GetPlayerInfo) Handle(player *Player) {
}
