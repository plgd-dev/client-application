// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: github.com/plgd-dev/client-application/pb/get_devices.proto

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

type GetDevicesRequest_OwnershipStatusFilter int32

const (
	// get only not owned devices
	GetDevicesRequest_UNOWNED GetDevicesRequest_OwnershipStatusFilter = 0
	// get only owned devices
	GetDevicesRequest_OWNED GetDevicesRequest_OwnershipStatusFilter = 1
)

// Enum value maps for GetDevicesRequest_OwnershipStatusFilter.
var (
	GetDevicesRequest_OwnershipStatusFilter_name = map[int32]string{
		0: "UNOWNED",
		1: "OWNED",
	}
	GetDevicesRequest_OwnershipStatusFilter_value = map[string]int32{
		"UNOWNED": 0,
		"OWNED":   1,
	}
)

func (x GetDevicesRequest_OwnershipStatusFilter) Enum() *GetDevicesRequest_OwnershipStatusFilter {
	p := new(GetDevicesRequest_OwnershipStatusFilter)
	*p = x
	return p
}

func (x GetDevicesRequest_OwnershipStatusFilter) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GetDevicesRequest_OwnershipStatusFilter) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes[0].Descriptor()
}

func (GetDevicesRequest_OwnershipStatusFilter) Type() protoreflect.EnumType {
	return &file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes[0]
}

func (x GetDevicesRequest_OwnershipStatusFilter) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GetDevicesRequest_OwnershipStatusFilter.Descriptor instead.
func (GetDevicesRequest_OwnershipStatusFilter) EnumDescriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescGZIP(), []int{0, 0}
}

type GetDevicesRequest_IPAddressType int32

const (
	GetDevicesRequest_IPV4 GetDevicesRequest_IPAddressType = 0
	GetDevicesRequest_IPV6 GetDevicesRequest_IPAddressType = 1
)

// Enum value maps for GetDevicesRequest_IPAddressType.
var (
	GetDevicesRequest_IPAddressType_name = map[int32]string{
		0: "IPV4",
		1: "IPV6",
	}
	GetDevicesRequest_IPAddressType_value = map[string]int32{
		"IPV4": 0,
		"IPV6": 1,
	}
)

func (x GetDevicesRequest_IPAddressType) Enum() *GetDevicesRequest_IPAddressType {
	p := new(GetDevicesRequest_IPAddressType)
	*p = x
	return p
}

func (x GetDevicesRequest_IPAddressType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GetDevicesRequest_IPAddressType) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes[1].Descriptor()
}

func (GetDevicesRequest_IPAddressType) Type() protoreflect.EnumType {
	return &file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes[1]
}

func (x GetDevicesRequest_IPAddressType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GetDevicesRequest_IPAddressType.Descriptor instead.
func (GetDevicesRequest_IPAddressType) EnumDescriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescGZIP(), []int{0, 1}
}

type GetDevicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Filter by ownership status. Default: [UNOWNED,OWNED].
	OwnershipStatusFilter []GetDevicesRequest_OwnershipStatusFilter `protobuf:"varint,1,rep,packed,name=ownership_status_filter,json=ownershipStatusFilter,proto3,enum=service.pb.GetDevicesRequest_OwnershipStatusFilter" json:"ownership_status_filter,omitempty"`
	// Filter by multicast IP address type. Default: [IPV4,IPV6].
	IpTypeFilter []GetDevicesRequest_IPAddressType `protobuf:"varint,2,rep,packed,name=ip_type_filter,json=ipTypeFilter,proto3,enum=service.pb.GetDevicesRequest_IPAddressType" json:"ip_type_filter,omitempty"`
	// Filter by device resource type of oic/d. Default: [] - filter is disabled.
	TypeFilter []string `protobuf:"bytes,4,rep,name=type_filter,json=typeFilter,proto3" json:"type_filter,omitempty"`
	// How long to wait for the devices responses for multicast request. 0 means use cache. Default: 0.
	Discover int64 `protobuf:"varint,5,opt,name=discover,proto3" json:"discover,omitempty"`
}

func (x *GetDevicesRequest) Reset() {
	*x = GetDevicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_get_devices_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDevicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDevicesRequest) ProtoMessage() {}

func (x *GetDevicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_get_devices_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDevicesRequest.ProtoReflect.Descriptor instead.
func (*GetDevicesRequest) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescGZIP(), []int{0}
}

func (x *GetDevicesRequest) GetOwnershipStatusFilter() []GetDevicesRequest_OwnershipStatusFilter {
	if x != nil {
		return x.OwnershipStatusFilter
	}
	return nil
}

func (x *GetDevicesRequest) GetIpTypeFilter() []GetDevicesRequest_IPAddressType {
	if x != nil {
		return x.IpTypeFilter
	}
	return nil
}

func (x *GetDevicesRequest) GetTypeFilter() []string {
	if x != nil {
		return x.TypeFilter
	}
	return nil
}

func (x *GetDevicesRequest) GetDiscover() int64 {
	if x != nil {
		return x.Discover
	}
	return 0
}

var File_github_com_plgd_dev_client_application_pb_get_devices_proto protoreflect.FileDescriptor

var file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDesc = []byte{
	0x0a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67,
	0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x74, 0x5f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x22, 0xe6, 0x02, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x6b, 0x0a, 0x17, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x33, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65,
	0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x4f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x15, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x51, 0x0a, 0x0e,
	0x69, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x62, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0c, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x79, 0x70, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x22, 0x2f, 0x0a, 0x15,
	0x4f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4f, 0x57, 0x4e, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x57, 0x4e, 0x45, 0x44, 0x10, 0x01, 0x22, 0x23, 0x0a,
	0x0d, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x49, 0x50, 0x56, 0x34, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x50, 0x56, 0x36,
	0x10, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x6c, 0x67, 0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x3b,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescOnce sync.Once
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescData = file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDesc
)

func file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescGZIP() []byte {
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescOnce.Do(func() {
		file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescData)
	})
	return file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDescData
}

var file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_github_com_plgd_dev_client_application_pb_get_devices_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_plgd_dev_client_application_pb_get_devices_proto_goTypes = []interface{}{
	(GetDevicesRequest_OwnershipStatusFilter)(0), // 0: service.pb.GetDevicesRequest.OwnershipStatusFilter
	(GetDevicesRequest_IPAddressType)(0),         // 1: service.pb.GetDevicesRequest.IPAddressType
	(*GetDevicesRequest)(nil),                    // 2: service.pb.GetDevicesRequest
}
var file_github_com_plgd_dev_client_application_pb_get_devices_proto_depIdxs = []int32{
	0, // 0: service.pb.GetDevicesRequest.ownership_status_filter:type_name -> service.pb.GetDevicesRequest.OwnershipStatusFilter
	1, // 1: service.pb.GetDevicesRequest.ip_type_filter:type_name -> service.pb.GetDevicesRequest.IPAddressType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_github_com_plgd_dev_client_application_pb_get_devices_proto_init() }
func file_github_com_plgd_dev_client_application_pb_get_devices_proto_init() {
	if File_github_com_plgd_dev_client_application_pb_get_devices_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_plgd_dev_client_application_pb_get_devices_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDevicesRequest); i {
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
			RawDescriptor: file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_plgd_dev_client_application_pb_get_devices_proto_goTypes,
		DependencyIndexes: file_github_com_plgd_dev_client_application_pb_get_devices_proto_depIdxs,
		EnumInfos:         file_github_com_plgd_dev_client_application_pb_get_devices_proto_enumTypes,
		MessageInfos:      file_github_com_plgd_dev_client_application_pb_get_devices_proto_msgTypes,
	}.Build()
	File_github_com_plgd_dev_client_application_pb_get_devices_proto = out.File
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_rawDesc = nil
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_goTypes = nil
	file_github_com_plgd_dev_client_application_pb_get_devices_proto_depIdxs = nil
}
