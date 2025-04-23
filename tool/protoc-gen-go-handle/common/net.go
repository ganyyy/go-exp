package common

import "protoc-gen-go-handle/protogen"

type INetServer interface {
	HandleNetMsg(netReq *protogen.NetReq) (*protogen.NetRsp, error)
}
