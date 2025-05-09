// Code generated by protoc-gen-go-handle. DO NOT EDIT.
package battlex

import (
	context "context"
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	common "protoc-gen-go-handle/common"
	protogen "protoc-gen-go-handle/protogen"
)

var (
	_ = fmt.Println
	_ = context.Background()
	_ = proto.Marshal
	_ = common.Version
	_ = (*protogen.Empty)(nil)
)

const (
	RoomMgrCreateRoom        = "RoomMgr.CreateRoom"
	RoomMgrJoinRoom          = "RoomMgr.JoinRoom"
	RoomMgrLeaveRoom         = "RoomMgr.LeaveRoom"
	RoomMgrNotifyBattleStart = "RoomMgr.NotifyBattleStart"
)

type innerRoomMgr struct {
	impl RoomMgrImpl
	ch   common.ChannelInvoke
}

type RoomMgr struct {
	inner innerRoomMgr
}

func (m *innerRoomMgr) HandleNetMsg(netReq *protogen.NetReq) (r *protogen.NetRsp, _ error) {
	var netRsp protogen.NetRsp
	var req, resp proto.Message
	var isAsync bool
	switch netReq.Code {
	case RoomMgrCreateRoom:
		req = new(protogen.CreateRoomReq)
		resp = new(protogen.CreateRoomRsp)
	case RoomMgrJoinRoom:
		req = new(protogen.JoinRoomReq)
		resp = new(protogen.JoinRoomRsp)
	case RoomMgrLeaveRoom:
		req = new(protogen.LeaveRoomReq)
		resp = new(protogen.LeaveRoomRsp)
	case RoomMgrNotifyBattleStart:
		req = new(protogen.NotifyBattleStartReq)
		isAsync = true
	default:
		return nil, fmt.Errorf("unknown task code: %v", netReq.Code)
	}
	if err := proto.Unmarshal(netReq.Req, req); err != nil {
		return nil, err
	}
	ctx, cancel := common.GenerateOptions().Context()
	defer cancel()
	err := m.ch.Invoke(ctx, netReq.Code, req, resp)
	if isAsync {
		return nil, err
	}
	if err != nil {
		netRsp.Err = err.Error()
	} else {
		netRsp.Rsp, _ = proto.Marshal(resp)
	}
	return &netRsp, err
}

func (m *innerRoomMgr) HandleTask(task common.Task) {
	defer func() {
		if r := recover(); r != nil {
			task.Finish(fmt.Errorf("task panic: %v", r))
		}
	}()
	switch task.Code {
	case RoomMgrCreateRoom:
		req := task.Req.(*protogen.CreateRoomReq)
		resp := task.Rsp.(*protogen.CreateRoomRsp)
		task.Finish(m.impl.CreateRoom(task.Context, req, resp))
	case RoomMgrJoinRoom:
		req := task.Req.(*protogen.JoinRoomReq)
		resp := task.Rsp.(*protogen.JoinRoomRsp)
		task.Finish(m.impl.JoinRoom(task.Context, req, resp))
	case RoomMgrLeaveRoom:
		req := task.Req.(*protogen.LeaveRoomReq)
		resp := task.Rsp.(*protogen.LeaveRoomRsp)
		task.Finish(m.impl.LeaveRoom(task.Context, req, resp))
	case RoomMgrNotifyBattleStart:
		req := task.Req.(*protogen.NotifyBattleStartReq)
		m.impl.NotifyBattleStart(task.Context, req)
		task.Finish(nil)
	default:
		task.Finish(fmt.Errorf("unknown task code: %v", task.Code))
	}
}

func (m *RoomMgr) L() IRoomMgrClient {
	return NewRoomMgrClient(&m.inner.ch)
}

type IRoomMgrClient interface {
	CreateRoom(*protogen.CreateRoomReq, ...common.ApplyOption) (*protogen.CreateRoomRsp, error)
	JoinRoom(*protogen.JoinRoomReq, ...common.ApplyOption) (*protogen.JoinRoomRsp, error)
	LeaveRoom(*protogen.LeaveRoomReq, ...common.ApplyOption) (*protogen.LeaveRoomRsp, error)
	NotifyBattleStart(*protogen.NotifyBattleStartReq, ...common.ApplyOption) error
}

type iRoomMgrClient struct {
	cc common.ClientInterface
}

func NewRoomMgrClient(cc common.ClientInterface) IRoomMgrClient {
	return iRoomMgrClient{cc}
}

func (c iRoomMgrClient) CreateRoom(in *protogen.CreateRoomReq, opts ...common.ApplyOption) (*protogen.CreateRoomRsp, error) {
	ctx, cancel := common.GenerateOptions(opts...).Context()
	defer cancel()
	out := new(protogen.CreateRoomRsp)
	err := c.cc.Invoke(ctx, RoomMgrCreateRoom, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c iRoomMgrClient) JoinRoom(in *protogen.JoinRoomReq, opts ...common.ApplyOption) (*protogen.JoinRoomRsp, error) {
	ctx, cancel := common.GenerateOptions(opts...).Context()
	defer cancel()
	out := new(protogen.JoinRoomRsp)
	err := c.cc.Invoke(ctx, RoomMgrJoinRoom, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c iRoomMgrClient) LeaveRoom(in *protogen.LeaveRoomReq, opts ...common.ApplyOption) (*protogen.LeaveRoomRsp, error) {
	ctx, cancel := common.GenerateOptions(opts...).Context()
	defer cancel()
	out := new(protogen.LeaveRoomRsp)
	err := c.cc.Invoke(ctx, RoomMgrLeaveRoom, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c iRoomMgrClient) NotifyBattleStart(in *protogen.NotifyBattleStartReq, opts ...common.ApplyOption) error {
	ctx, cancel := common.GenerateOptions(opts...).Context()
	defer cancel()
	return c.cc.Invoke(ctx, RoomMgrNotifyBattleStart, in, nil)
}
