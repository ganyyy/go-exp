// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.17.3
// source: proto/player/player.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	data "protoc-gen-readonly/pb/data"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type State int32

const (
	State_Unknown State = 0
	State_Login   State = 1
	State_Logout  State = 2
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "Unknown",
		1: "Login",
		2: "Logout",
	}
	State_value = map[string]int32{
		"Unknown": 0,
		"Login":   1,
		"Logout":  2,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_player_player_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_proto_player_player_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_proto_player_player_proto_rawDescGZIP(), []int{0}
}

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string                         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Sa     []*data.SimpleData             `protobuf:"bytes,2,rep,name=sa,proto3" json:"sa,omitempty"`
	Ma     map[int32]*data.ReferencedData `protobuf:"bytes,3,rep,name=ma,proto3" json:"ma,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Da     *data.SimpleData               `protobuf:"bytes,4,opt,name=da,proto3" json:"da,omitempty"`
	Age    *int32                         `protobuf:"varint,5,opt,name=age,proto3,oneof" json:"age,omitempty"`
	State  State                          `protobuf:"varint,6,opt,name=state,proto3,enum=player.State" json:"state,omitempty"`
	State2 *State                         `protobuf:"varint,7,opt,name=state2,proto3,enum=player.State,oneof" json:"state2,omitempty"`
	Data   []byte                         `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	mi := &file_proto_player_player_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_proto_player_player_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_proto_player_player_proto_rawDescGZIP(), []int{0}
}

func (x *Player) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Player) GetSa() []*data.SimpleData {
	if x != nil {
		return x.Sa
	}
	return nil
}

func (x *Player) GetMa() map[int32]*data.ReferencedData {
	if x != nil {
		return x.Ma
	}
	return nil
}

func (x *Player) GetDa() *data.SimpleData {
	if x != nil {
		return x.Da
	}
	return nil
}

func (x *Player) GetAge() int32 {
	if x != nil && x.Age != nil {
		return *x.Age
	}
	return 0
}

func (x *Player) GetState() State {
	if x != nil {
		return x.State
	}
	return State_Unknown
}

func (x *Player) GetState2() State {
	if x != nil && x.State2 != nil {
		return *x.State2
	}
	return State_Unknown
}

func (x *Player) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_player_player_proto protoreflect.FileDescriptor

var file_proto_player_player_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2f, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x1a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x02, 0x0a, 0x06, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x02, 0x73, 0x61, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x02, 0x73, 0x61, 0x12, 0x26, 0x0a, 0x02, 0x6d,
	0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x02, 0x6d, 0x61, 0x12, 0x20, 0x0a, 0x02, 0x64, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x02, 0x64, 0x61, 0x12, 0x15, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x00, 0x52, 0x03, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x65, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0d, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x48, 0x01, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x65, 0x32, 0x88, 0x01, 0x01, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x1a, 0x4b, 0x0a, 0x07, 0x4d, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x06,
	0x0a, 0x04, 0x5f, 0x61, 0x67, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x32, 0x2a, 0x2b, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e,
	0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x10, 0x02, 0x42, 0x18,
	0x5a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x72, 0x65, 0x61,
	0x64, 0x6f, 0x6e, 0x6c, 0x79, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_player_player_proto_rawDescOnce sync.Once
	file_proto_player_player_proto_rawDescData = file_proto_player_player_proto_rawDesc
)

func file_proto_player_player_proto_rawDescGZIP() []byte {
	file_proto_player_player_proto_rawDescOnce.Do(func() {
		file_proto_player_player_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_player_player_proto_rawDescData)
	})
	return file_proto_player_player_proto_rawDescData
}

var file_proto_player_player_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_player_player_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_player_player_proto_goTypes = []any{
	(State)(0),                  // 0: player.State
	(*Player)(nil),              // 1: player.Player
	nil,                         // 2: player.Player.MaEntry
	(*data.SimpleData)(nil),     // 3: data.SimpleData
	(*data.ReferencedData)(nil), // 4: data.ReferencedData
}
var file_proto_player_player_proto_depIdxs = []int32{
	3, // 0: player.Player.sa:type_name -> data.SimpleData
	2, // 1: player.Player.ma:type_name -> player.Player.MaEntry
	3, // 2: player.Player.da:type_name -> data.SimpleData
	0, // 3: player.Player.state:type_name -> player.State
	0, // 4: player.Player.state2:type_name -> player.State
	4, // 5: player.Player.MaEntry.value:type_name -> data.ReferencedData
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_player_player_proto_init() }
func file_proto_player_player_proto_init() {
	if File_proto_player_player_proto != nil {
		return
	}
	file_proto_player_player_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_player_player_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_player_player_proto_goTypes,
		DependencyIndexes: file_proto_player_player_proto_depIdxs,
		EnumInfos:         file_proto_player_player_proto_enumTypes,
		MessageInfos:      file_proto_player_player_proto_msgTypes,
	}.Build()
	File_proto_player_player_proto = out.File
	file_proto_player_player_proto_rawDesc = nil
	file_proto_player_player_proto_goTypes = nil
	file_proto_player_player_proto_depIdxs = nil
}