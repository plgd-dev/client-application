// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: github.com/plgd-dev/client-application/pb/reset.proto

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

type ResetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ResetRequest) Reset() {
	*x = ResetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetRequest) ProtoMessage() {}

func (x *ResetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetRequest.ProtoReflect.Descriptor instead.
func (*ResetRequest) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescGZIP(), []int{0}
}

type ResetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ResetResponse) Reset() {
	*x = ResetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetResponse) ProtoMessage() {}

func (x *ResetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetResponse.ProtoReflect.Descriptor instead.
func (*ResetResponse) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescGZIP(), []int{1}
}

var File_github_com_plgd_dev_client_application_pb_reset_proto protoreflect.FileDescriptor

var file_github_com_plgd_dev_client_application_pb_reset_proto_rawDesc = []byte{
	0x0a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67,
	0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x73, 0x65,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x62, 0x22, 0x0e, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x0f, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67, 0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70,
	0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescOnce sync.Once
	file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescData = file_github_com_plgd_dev_client_application_pb_reset_proto_rawDesc
)

func file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescGZIP() []byte {
	file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescOnce.Do(func() {
		file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescData)
	})
	return file_github_com_plgd_dev_client_application_pb_reset_proto_rawDescData
}

var file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_plgd_dev_client_application_pb_reset_proto_goTypes = []interface{}{
	(*ResetRequest)(nil),  // 0: service.pb.ResetRequest
	(*ResetResponse)(nil), // 1: service.pb.ResetResponse
}
var file_github_com_plgd_dev_client_application_pb_reset_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_plgd_dev_client_application_pb_reset_proto_init() }
func file_github_com_plgd_dev_client_application_pb_reset_proto_init() {
	if File_github_com_plgd_dev_client_application_pb_reset_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResetRequest); i {
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
		file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResetResponse); i {
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
			RawDescriptor: file_github_com_plgd_dev_client_application_pb_reset_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_plgd_dev_client_application_pb_reset_proto_goTypes,
		DependencyIndexes: file_github_com_plgd_dev_client_application_pb_reset_proto_depIdxs,
		MessageInfos:      file_github_com_plgd_dev_client_application_pb_reset_proto_msgTypes,
	}.Build()
	File_github_com_plgd_dev_client_application_pb_reset_proto = out.File
	file_github_com_plgd_dev_client_application_pb_reset_proto_rawDesc = nil
	file_github_com_plgd_dev_client_application_pb_reset_proto_goTypes = nil
	file_github_com_plgd_dev_client_application_pb_reset_proto_depIdxs = nil
}
