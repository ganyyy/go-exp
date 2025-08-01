// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.7
// source: battle.proto

package protogen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Battle struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BattleId      string                 `protobuf:"bytes,1,opt,name=battleId,proto3" json:"battleId,omitempty"`
	PlayerIds     []string               `protobuf:"bytes,2,rep,name=playerIds,proto3" json:"playerIds,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Battle) Reset() {
	*x = Battle{}
	mi := &file_battle_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Battle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Battle) ProtoMessage() {}

func (x *Battle) ProtoReflect() protoreflect.Message {
	mi := &file_battle_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Battle.ProtoReflect.Descriptor instead.
func (*Battle) Descriptor() ([]byte, []int) {
	return file_battle_proto_rawDescGZIP(), []int{0}
}

func (x *Battle) GetBattleId() string {
	if x != nil {
		return x.BattleId
	}
	return ""
}

func (x *Battle) GetPlayerIds() []string {
	if x != nil {
		return x.PlayerIds
	}
	return nil
}

type GetBattleInfoReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BattleId      string                 `protobuf:"bytes,1,opt,name=battleId,proto3" json:"battleId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBattleInfoReq) Reset() {
	*x = GetBattleInfoReq{}
	mi := &file_battle_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBattleInfoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBattleInfoReq) ProtoMessage() {}

func (x *GetBattleInfoReq) ProtoReflect() protoreflect.Message {
	mi := &file_battle_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBattleInfoReq.ProtoReflect.Descriptor instead.
func (*GetBattleInfoReq) Descriptor() ([]byte, []int) {
	return file_battle_proto_rawDescGZIP(), []int{1}
}

func (x *GetBattleInfoReq) GetBattleId() string {
	if x != nil {
		return x.BattleId
	}
	return ""
}

type GetBattleInfoRsp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Battle        *Battle                `protobuf:"bytes,1,opt,name=battle,proto3" json:"battle,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBattleInfoRsp) Reset() {
	*x = GetBattleInfoRsp{}
	mi := &file_battle_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBattleInfoRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBattleInfoRsp) ProtoMessage() {}

func (x *GetBattleInfoRsp) ProtoReflect() protoreflect.Message {
	mi := &file_battle_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBattleInfoRsp.ProtoReflect.Descriptor instead.
func (*GetBattleInfoRsp) Descriptor() ([]byte, []int) {
	return file_battle_proto_rawDescGZIP(), []int{2}
}

func (x *GetBattleInfoRsp) GetBattle() *Battle {
	if x != nil {
		return x.Battle
	}
	return nil
}

var File_battle_proto protoreflect.FileDescriptor

const file_battle_proto_rawDesc = "" +
	"\n" +
	"\fbattle.proto\x12\abattlex\"B\n" +
	"\x06Battle\x12\x1a\n" +
	"\bbattleId\x18\x01 \x01(\tR\bbattleId\x12\x1c\n" +
	"\tplayerIds\x18\x02 \x03(\tR\tplayerIds\".\n" +
	"\x10GetBattleInfoReq\x12\x1a\n" +
	"\bbattleId\x18\x01 \x01(\tR\bbattleId\";\n" +
	"\x10GetBattleInfoRsp\x12'\n" +
	"\x06battle\x18\x01 \x01(\v2\x0f.battlex.BattleR\x06battle2\x98\x01\n" +
	"\aBattleX\x12E\n" +
	"\rGetBattleInfo\x12\x19.battlex.GetBattleInfoReq\x1a\x19.battlex.GetBattleInfoRsp\x12F\n" +
	"\x0eGetBattleInfo2\x12\x19.battlex.GetBattleInfoReq\x1a\x19.battlex.GetBattleInfoRspB\x1fZ\x1dprotoc-gen-go-handle/protogenb\x06proto3"

var (
	file_battle_proto_rawDescOnce sync.Once
	file_battle_proto_rawDescData []byte
)

func file_battle_proto_rawDescGZIP() []byte {
	file_battle_proto_rawDescOnce.Do(func() {
		file_battle_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_battle_proto_rawDesc), len(file_battle_proto_rawDesc)))
	})
	return file_battle_proto_rawDescData
}

var file_battle_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_battle_proto_goTypes = []any{
	(*Battle)(nil),           // 0: battlex.Battle
	(*GetBattleInfoReq)(nil), // 1: battlex.GetBattleInfoReq
	(*GetBattleInfoRsp)(nil), // 2: battlex.GetBattleInfoRsp
}
var file_battle_proto_depIdxs = []int32{
	0, // 0: battlex.GetBattleInfoRsp.battle:type_name -> battlex.Battle
	1, // 1: battlex.BattleX.GetBattleInfo:input_type -> battlex.GetBattleInfoReq
	1, // 2: battlex.BattleX.GetBattleInfo2:input_type -> battlex.GetBattleInfoReq
	2, // 3: battlex.BattleX.GetBattleInfo:output_type -> battlex.GetBattleInfoRsp
	2, // 4: battlex.BattleX.GetBattleInfo2:output_type -> battlex.GetBattleInfoRsp
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_battle_proto_init() }
func file_battle_proto_init() {
	if File_battle_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_battle_proto_rawDesc), len(file_battle_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_battle_proto_goTypes,
		DependencyIndexes: file_battle_proto_depIdxs,
		MessageInfos:      file_battle_proto_msgTypes,
	}.Build()
	File_battle_proto = out.File
	file_battle_proto_goTypes = nil
	file_battle_proto_depIdxs = nil
}
