package battlex

import (
	context "context"
	protogen "protoc-gen-go-handle/protogen"
)

var (
	_ = context.Background()
	_ = (*protogen.Empty)(nil)
)

type RoomMgrImpl struct {
}

func (s *RoomMgrImpl) CreateRoom(ctx context.Context, req *protogen.CreateRoomReq, rsp *protogen.CreateRoomRsp) error {
	return nil
}

func (s *RoomMgrImpl) JoinRoom(ctx context.Context, req *protogen.JoinRoomReq, rsp *protogen.JoinRoomRsp) error {
	return nil
}

func (s *RoomMgrImpl) LeaveRoom(ctx context.Context, req *protogen.LeaveRoomReq, rsp *protogen.LeaveRoomRsp) error {
	return nil
}

func (s *RoomMgrImpl) NotifyBattleStart(ctx context.Context, req *protogen.NotifyBattleStartReq) {
}
