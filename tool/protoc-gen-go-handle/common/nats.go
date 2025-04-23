package common

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"protoc-gen-go-handle/protogen"
)

type NatsInvoke string

// Invoke 通过nats调用
func (n NatsInvoke) Invoke(ctx context.Context, code string, req proto.Message, rsp proto.Message) error {
	t, ok := ctx.Deadline()
	var timeout time.Duration
	if ok {
		timeout = time.Until(t)
	} else {
		timeout = time.Second
	}
	_ = timeout
	var netReq protogen.NetReq
	netReq.Code = code
	netReq.Req, _ = proto.Marshal(req)

	bs, _ := proto.Marshal(&netReq)

	if rsp == nil {
		// TODO 异步任务, Publish
		return nil
	}
	var netResp protogen.NetRsp
	// TODO 同步任务, Request
	_ = bs

	var respBytes []byte
	_ = proto.Unmarshal(respBytes, &netResp)
	if netResp.Err != "" {
		return errors.New(netResp.Err)
	}
	_ = proto.Unmarshal(netResp.Rsp, rsp)

	return nil

}
