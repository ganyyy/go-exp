package msg

import (
	"testing"

	pb "ganyyy.com/go-exp/rpc/grpc/proto"
	"github.com/golang/protobuf/proto"
)

func TestProtoMarshal(t *testing.T) {
	var msg = pb.FieldType{
		F1: -1,
	}

	var size = proto.Size(&msg)

	t.Logf("size:%v", size)

	var data, _ = proto.Marshal(&msg)

	t.Logf("%v", string(data))
}
