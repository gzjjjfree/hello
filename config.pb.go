// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0--rc2
// source: config.proto

package core

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Inbound handler configurations. Must have at least one item.
	// 入站处理程序配置。必须至少有一个项目
	Inbounds []*InboundHandlerConfig `protobuf:"bytes,1,rep,name=inbounds,proto3" json:"inbounds,omitempty"`
	// Outbound handler configurations. Must have at least one item. The first
	// item is used as default for routing.
	// 出站处理程序配置。必须至少有一个项目。第一个项目用作路由的默认值。
	Outbounds []*OutboundHandlerConfig `protobuf:"bytes,2,rep,name=outbounds,proto3" json:"outbounds,omitempty"`
	// App is for configurations of all features in V2Ray. A feature must
	// App 用于配置 V2Ray 中的所有功能。功能必须
	// implement the Feature interface, and its config type must be registered
	// 实现 Feature 接口，并且必须注册其配置类型，通过common.RegisterConfig。
	// through common.RegisterConfig.
	App []*TypedMessage `protobuf:"bytes,4,rep,name=app,proto3" json:"app,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetInbounds() []*InboundHandlerConfig {
	if x != nil {
		return x.Inbounds
	}
	return nil
}

func (x *Config) GetOutbounds() []*OutboundHandlerConfig {
	if x != nil {
		return x.Outbounds
	}
	return nil
}

func (x *Config) GetApp() []*TypedMessage {
	if x != nil {
		return x.App
	}
	return nil
}

// InboundHandlerConfig is the configuration for inbound handler.
// InboundHandlerConfig 是入站处理程序的配置。
type InboundHandlerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag      string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Protocol string `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *InboundHandlerConfig) Reset() {
	*x = InboundHandlerConfig{}
	mi := &file_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InboundHandlerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InboundHandlerConfig) ProtoMessage() {}

func (x *InboundHandlerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InboundHandlerConfig.ProtoReflect.Descriptor instead.
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{1}
}

func (x *InboundHandlerConfig) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *InboundHandlerConfig) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// OutboundHandlerConfig is the configuration for outbound handler.
// OutboundHandlerConfig 是出站处理程序的配置。
type OutboundHandlerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag      string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Protocol string `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *OutboundHandlerConfig) Reset() {
	*x = OutboundHandlerConfig{}
	mi := &file_config_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OutboundHandlerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutboundHandlerConfig) ProtoMessage() {}

func (x *OutboundHandlerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutboundHandlerConfig.ProtoReflect.Descriptor instead.
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *OutboundHandlerConfig) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *OutboundHandlerConfig) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// TypedMessage is a serialized proto message along with its type name.
type TypedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the message type, retrieved from protobuf API.
	// 从 protobuf API 中检索的消息类型的名称
	// string type = 1;
	// Serialized proto message. 序列化的原始消息。
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *TypedMessage) Reset() {
	*x = TypedMessage{}
	mi := &file_config_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TypedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypedMessage) ProtoMessage() {}

func (x *TypedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypedMessage.ProtoReflect.Descriptor instead.
func (*TypedMessage) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
}

func (x *TypedMessage) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_config_proto protoreflect.FileDescriptor

var file_config_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x22, 0xb9, 0x01, 0x0a, 0x06, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3c, 0x0a, 0x08, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x08, 0x69, 0x6e, 0x62, 0x6f, 0x75,
	0x6e, 0x64, 0x73, 0x12, 0x3f, 0x0a, 0x09, 0x6f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x6f, 0x75, 0x74, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x73, 0x12, 0x2a, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x54,
	0x79, 0x70, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x61, 0x70, 0x70,
	0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x44, 0x0a, 0x14, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e,
	0x64, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x10,
	0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x45, 0x0a, 0x15,
	0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x22, 0x24, 0x0a, 0x0c, 0x54, 0x79, 0x70, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x7a, 0x6a, 0x6a, 0x6a, 0x66, 0x72, 0x65,
	0x65, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x3b, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_proto_rawDescOnce sync.Once
	file_config_proto_rawDescData = file_config_proto_rawDesc
)

func file_config_proto_rawDescGZIP() []byte {
	file_config_proto_rawDescOnce.Do(func() {
		file_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_proto_rawDescData)
	})
	return file_config_proto_rawDescData
}

var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_proto_goTypes = []any{
	(*Config)(nil),                // 0: hello.core.Config
	(*InboundHandlerConfig)(nil),  // 1: hello.core.InboundHandlerConfig
	(*OutboundHandlerConfig)(nil), // 2: hello.core.OutboundHandlerConfig
	(*TypedMessage)(nil),          // 3: hello.core.TypedMessage
}
var file_config_proto_depIdxs = []int32{
	1, // 0: hello.core.Config.inbounds:type_name -> hello.core.InboundHandlerConfig
	2, // 1: hello.core.Config.outbounds:type_name -> hello.core.OutboundHandlerConfig
	3, // 2: hello.core.Config.app:type_name -> hello.core.TypedMessage
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_config_proto_init() }
func file_config_proto_init() {
	if File_config_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_rawDesc = nil
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}