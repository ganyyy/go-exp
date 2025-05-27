package battlex

import (
	context "context"
	protogen "protoc-gen-go-handle/protogen"
)

var (
	_ = context.Background()
	_ = (*protogen.Empty)(nil)
)

type BattleXImpl struct {
}

// GetBattleInfo2
func (s *BattleXImpl) GetBattleInfo2(ctx context.Context, req *protogen.GetBattleInfoReq, rsp *protogen.GetBattleInfoRsp) error {
	return nil
}

// GetBattleInfo
func (s *BattleXImpl) GetBattleInfo(ctx context.Context, req *protogen.GetBattleInfoReq, rsp *protogen.GetBattleInfoRsp) error {
	return nil
}
