// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: github.com/plgd-dev/client-application/pb/get_resource.proto

package pb

import (
	commands "github.com/plgd-dev/hub/v2/resource-aggregate/commands"
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

type GetResourceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId        *commands.ResourceId `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ResourceInterface string               `protobuf:"bytes,2,opt,name=resource_interface,json=resourceInterface,proto3" json:"resource_interface,omitempty"`
}

func (x *GetResourceRequest) Reset() {
	*x = GetResourceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_get_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResourceRequest) ProtoMessage() {}

func (x *GetResourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_get_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResourceRequest.ProtoReflect.Descriptor instead.
func (*GetResourceRequest) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescGZIP(), []int{0}
}

func (x *GetResourceRequest) GetResourceId() *commands.ResourceId {
	if x != nil {
		return x.ResourceId
	}
	return nil
}

func (x *GetResourceRequest) GetResourceInterface() string {
	if x != nil {
		return x.ResourceInterface
	}
	return ""
}

var File_github_com_plgd_dev_client_application_pb_get_resource_proto protoreflect.FileDescriptor

var file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67,
	0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x74, 0x5f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x1a, 0x25, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2d, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x2f, 0x70,
	0x62, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x86, 0x01, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x52,
	0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x12, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67, 0x64, 0x2d, 0x64, 0x65,
	0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescOnce sync.Once
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescData = file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDesc
)

func file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescGZIP() []byte {
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescOnce.Do(func() {
		file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescData)
	})
	return file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDescData
}

var file_github_com_plgd_dev_client_application_pb_get_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_plgd_dev_client_application_pb_get_resource_proto_goTypes = []any{
	(*GetResourceRequest)(nil),  // 0: service.pb.GetResourceRequest
	(*commands.ResourceId)(nil), // 1: resourceaggregate.pb.ResourceId
}
var file_github_com_plgd_dev_client_application_pb_get_resource_proto_depIdxs = []int32{
	1, // 0: service.pb.GetResourceRequest.resource_id:type_name -> resourceaggregate.pb.ResourceId
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_plgd_dev_client_application_pb_get_resource_proto_init() }
func file_github_com_plgd_dev_client_application_pb_get_resource_proto_init() {
	if File_github_com_plgd_dev_client_application_pb_get_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_plgd_dev_client_application_pb_get_resource_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetResourceRequest); i {
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
			RawDescriptor: file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_plgd_dev_client_application_pb_get_resource_proto_goTypes,
		DependencyIndexes: file_github_com_plgd_dev_client_application_pb_get_resource_proto_depIdxs,
		MessageInfos:      file_github_com_plgd_dev_client_application_pb_get_resource_proto_msgTypes,
	}.Build()
	File_github_com_plgd_dev_client_application_pb_get_resource_proto = out.File
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_rawDesc = nil
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_goTypes = nil
	file_github_com_plgd_dev_client_application_pb_get_resource_proto_depIdxs = nil
}
