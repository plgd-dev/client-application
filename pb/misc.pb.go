// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: github.com/plgd-dev/client-application/pb/misc.proto

package pb

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

type Device_OwnershipStatus int32

const (
	// cannot determine ownership status
	Device_UNKNOWN Device_OwnershipStatus = 0
	// device is ready to be owned the user
	Device_UNOWNED Device_OwnershipStatus = 1
	// device is owned by the user. to determine who own the device you need to get ownership resource /oic/sec/doxm
	Device_OWNED Device_OwnershipStatus = 2
	// set when device is not secured. (iotivity-lite was built without security)
	Device_UNSUPPORTED Device_OwnershipStatus = 3
)

// Enum value maps for Device_OwnershipStatus.
var (
	Device_OwnershipStatus_name = map[int32]string{
		0: "UNKNOWN",
		1: "UNOWNED",
		2: "OWNED",
		3: "UNSUPPORTED",
	}
	Device_OwnershipStatus_value = map[string]int32{
		"UNKNOWN":     0,
		"UNOWNED":     1,
		"OWNED":       2,
		"UNSUPPORTED": 3,
	}
)

func (x Device_OwnershipStatus) Enum() *Device_OwnershipStatus {
	p := new(Device_OwnershipStatus)
	*p = x
	return p
}

func (x Device_OwnershipStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Device_OwnershipStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_plgd_dev_client_application_pb_misc_proto_enumTypes[0].Descriptor()
}

func (Device_OwnershipStatus) Type() protoreflect.EnumType {
	return &file_github_com_plgd_dev_client_application_pb_misc_proto_enumTypes[0]
}

func (x Device_OwnershipStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Device_OwnershipStatus.Descriptor instead.
func (Device_OwnershipStatus) EnumDescriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescGZIP(), []int{2, 0}
}

type ResourceId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId string `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	Href     string `protobuf:"bytes,2,opt,name=href,proto3" json:"href,omitempty"`
}

func (x *ResourceId) Reset() {
	*x = ResourceId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceId) ProtoMessage() {}

func (x *ResourceId) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceId.ProtoReflect.Descriptor instead.
func (*ResourceId) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceId) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *ResourceId) GetHref() string {
	if x != nil {
		return x.Href
	}
	return ""
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContentType string `protobuf:"bytes,1,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Data        []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescGZIP(), []int{1}
}

func (x *Content) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *Content) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// content of the device resource oic/d
	Content *Content `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	// ownership status of the device
	OwnershipStatus Device_OwnershipStatus `protobuf:"varint,2,opt,name=ownership_status,json=ownershipStatus,proto3,enum=service.pb.Device_OwnershipStatus" json:"ownership_status,omitempty"`
	// endpoints with schemas which are hosted by the device
	Endpoints []string `protobuf:"bytes,3,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescGZIP(), []int{2}
}

func (x *Device) GetContent() *Content {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *Device) GetOwnershipStatus() Device_OwnershipStatus {
	if x != nil {
		return x.OwnershipStatus
	}
	return Device_UNKNOWN
}

func (x *Device) GetEndpoints() []string {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

var File_github_com_plgd_dev_client_application_pb_misc_proto protoreflect.FileDescriptor

var file_github_com_plgd_dev_client_application_pb_misc_proto_rawDesc = []byte{
	0x0a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67,
	0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x6d, 0x69, 0x73, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x62, 0x22, 0x3d, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x72, 0x65, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x72, 0x65,
	0x66, 0x22, 0x40, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0xed, 0x01, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2d,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x4d, 0x0a,
	0x10, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4f, 0x77, 0x6e, 0x65,
	0x72, 0x73, 0x68, 0x69, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x6f, 0x77, 0x6e,
	0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x47, 0x0a, 0x0f, 0x4f, 0x77,
	0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4f, 0x57, 0x4e, 0x45, 0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x57, 0x4e, 0x45, 0x44,
	0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x55, 0x50, 0x50, 0x4f, 0x52, 0x54, 0x45,
	0x44, 0x10, 0x03, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x6c, 0x67, 0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62,
	0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescOnce sync.Once
	file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescData = file_github_com_plgd_dev_client_application_pb_misc_proto_rawDesc
)

func file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescGZIP() []byte {
	file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescOnce.Do(func() {
		file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescData)
	})
	return file_github_com_plgd_dev_client_application_pb_misc_proto_rawDescData
}

var file_github_com_plgd_dev_client_application_pb_misc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_github_com_plgd_dev_client_application_pb_misc_proto_goTypes = []interface{}{
	(Device_OwnershipStatus)(0), // 0: service.pb.Device.OwnershipStatus
	(*ResourceId)(nil),          // 1: service.pb.ResourceId
	(*Content)(nil),             // 2: service.pb.Content
	(*Device)(nil),              // 3: service.pb.Device
}
var file_github_com_plgd_dev_client_application_pb_misc_proto_depIdxs = []int32{
	2, // 0: service.pb.Device.content:type_name -> service.pb.Content
	0, // 1: service.pb.Device.ownership_status:type_name -> service.pb.Device.OwnershipStatus
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_github_com_plgd_dev_client_application_pb_misc_proto_init() }
func file_github_com_plgd_dev_client_application_pb_misc_proto_init() {
	if File_github_com_plgd_dev_client_application_pb_misc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_plgd_dev_client_application_pb_misc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_plgd_dev_client_application_pb_misc_proto_goTypes,
		DependencyIndexes: file_github_com_plgd_dev_client_application_pb_misc_proto_depIdxs,
		EnumInfos:         file_github_com_plgd_dev_client_application_pb_misc_proto_enumTypes,
		MessageInfos:      file_github_com_plgd_dev_client_application_pb_misc_proto_msgTypes,
	}.Build()
	File_github_com_plgd_dev_client_application_pb_misc_proto = out.File
	file_github_com_plgd_dev_client_application_pb_misc_proto_rawDesc = nil
	file_github_com_plgd_dev_client_application_pb_misc_proto_goTypes = nil
	file_github_com_plgd_dev_client_application_pb_misc_proto_depIdxs = nil
}