package player_module

import (
	context "context"
	"fmt"
	protogen "protoc-gen-go-handle/protogen"
)

var (
	_ = context.Background()
	_ = (*protogen.Empty)(nil)
	_ = fmt.Println
)

type Player struct {
	Acid string
	Name string
}
type PlayerModuleImpl struct {
	players map[string]*Player
}
type SomeTestStruct struct{}

func (s *PlayerModuleImpl) NotifyInfo(ctx context.Context, req *protogen.NotifyInfoReq) {
}
func (s *PlayerModuleImpl) GetPlayerInfo(ctx context.Context, req *protogen.GetPlayerInfoReq, rsp *protogen.GetPlayerInfoRsp) error {
	return nil
}
func (s *PlayerModuleImpl) GetName(ctx context.Context, req *protogen.GetNameReq, rsp *protogen.GetNameRsp) error {
	return nil
}
func (s *PlayerModuleImpl) NotifyAddAge(ctx context.Context, req *protogen.NotifyAddAgeReq) {
}
