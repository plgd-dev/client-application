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
// source: github.com/plgd-dev/client-application/pb/initialize.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InitializePreSharedKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the client application that is used to identify the client application by the device.
	SubjectUuid string `protobuf:"bytes,1,opt,name=subject_uuid,json=subjectUuid,proto3" json:"subject_uuid,omitempty"`
	// Associated secret to the client application ID.
	KeyUuid string `protobuf:"bytes,2,opt,name=key_uuid,json=keyUuid,proto3" json:"key_uuid,omitempty"`
}

func (x *InitializePreSharedKey) Reset() {
	*x = InitializePreSharedKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitializePreSharedKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializePreSharedKey) ProtoMessage() {}

func (x *InitializePreSharedKey) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializePreSharedKey.ProtoReflect.Descriptor instead.
func (*InitializePreSharedKey) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{0}
}

func (x *InitializePreSharedKey) GetSubjectUuid() string {
	if x != nil {
		return x.SubjectUuid
	}
	return ""
}

func (x *InitializePreSharedKey) GetKeyUuid() string {
	if x != nil {
		return x.KeyUuid
	}
	return ""
}

type InitializeX509 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Certificate chain in PEM format
	Certificate string `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"`
}

func (x *InitializeX509) Reset() {
	*x = InitializeX509{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitializeX509) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeX509) ProtoMessage() {}

func (x *InitializeX509) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeX509.ProtoReflect.Descriptor instead.
func (*InitializeX509) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{1}
}

func (x *InitializeX509) GetCertificate() string {
	if x != nil {
		return x.Certificate
	}
	return ""
}

// The client application must be initialized when GetConfigurationResponse.is_initialized is set to false.
// The initialization depends on the GetConfigurationResponse.device_authentication_mode.
// For:
//  - PRE_SHARED_KEY - pre_shared_key values need to be set.
//  - X509 - this values need to be set only if GetConfigurationResponse.remote_provisioning.mode is set to USER_AGENT, for SELF just use empty values
type InitializeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreSharedKey *InitializePreSharedKey `protobuf:"bytes,1,opt,name=pre_shared_key,json=preSharedKey,proto3" json:"pre_shared_key,omitempty"`
	Jwks         *structpb.Struct        `protobuf:"bytes,2,opt,name=jwks,proto3" json:"jwks,omitempty"`
}

func (x *InitializeRequest) Reset() {
	*x = InitializeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitializeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeRequest) ProtoMessage() {}

func (x *InitializeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeRequest.ProtoReflect.Descriptor instead.
func (*InitializeRequest) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{2}
}

func (x *InitializeRequest) GetPreSharedKey() *InitializePreSharedKey {
	if x != nil {
		return x.PreSharedKey
	}
	return nil
}

func (x *InitializeRequest) GetJwks() *structpb.Struct {
	if x != nil {
		return x.Jwks
	}
	return nil
}

type IdentityCertificateChallenge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CertificateSigningRequest []byte `protobuf:"bytes,1,opt,name=certificate_signing_request,json=certificateSigningRequest,proto3" json:"certificate_signing_request,omitempty"` // in PEM format
	State                     string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`                                                                            // for pairing calls
}

func (x *IdentityCertificateChallenge) Reset() {
	*x = IdentityCertificateChallenge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentityCertificateChallenge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityCertificateChallenge) ProtoMessage() {}

func (x *IdentityCertificateChallenge) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentityCertificateChallenge.ProtoReflect.Descriptor instead.
func (*IdentityCertificateChallenge) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{3}
}

func (x *IdentityCertificateChallenge) GetCertificateSigningRequest() []byte {
	if x != nil {
		return x.CertificateSigningRequest
	}
	return nil
}

func (x *IdentityCertificateChallenge) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type InitializeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If is set, the initialization process will be paused.
	// For the next call FinishInitialize, request must contain provided identity_certificate_challenge.state.
	IdentityCertificateChallenge *IdentityCertificateChallenge `protobuf:"bytes,2,opt,name=identity_certificate_challenge,json=identityCertificateChallenge,proto3" json:"identity_certificate_challenge,omitempty"`
}

func (x *InitializeResponse) Reset() {
	*x = InitializeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitializeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeResponse) ProtoMessage() {}

func (x *InitializeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeResponse.ProtoReflect.Descriptor instead.
func (*InitializeResponse) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{4}
}

func (x *InitializeResponse) GetIdentityCertificateChallenge() *IdentityCertificateChallenge {
	if x != nil {
		return x.IdentityCertificateChallenge
	}
	return nil
}

type FinishInitializeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Certificate chain in PEM format
	Certificate []byte `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"`
	// Use value for pairing otherwise finish will be refused.
	State string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *FinishInitializeRequest) Reset() {
	*x = FinishInitializeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishInitializeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishInitializeRequest) ProtoMessage() {}

func (x *FinishInitializeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishInitializeRequest.ProtoReflect.Descriptor instead.
func (*FinishInitializeRequest) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{5}
}

func (x *FinishInitializeRequest) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *FinishInitializeRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type FinishInitializeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FinishInitializeResponse) Reset() {
	*x = FinishInitializeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishInitializeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishInitializeResponse) ProtoMessage() {}

func (x *FinishInitializeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishInitializeResponse.ProtoReflect.Descriptor instead.
func (*FinishInitializeResponse) Descriptor() ([]byte, []int) {
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP(), []int{6}
}

var File_github_com_plgd_dev_client_application_pb_initialize_proto protoreflect.FileDescriptor

var file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDesc = []byte{
	0x0a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67,
	0x64, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x69, 0x6e, 0x69, 0x74,
	0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x16, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x50, 0x72, 0x65, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4b, 0x65, 0x79,
	0x12, 0x21, 0x0a, 0x0c, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x55,
	0x75, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x55, 0x75, 0x69, 0x64, 0x22, 0x32,
	0x0a, 0x0e, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x58, 0x35, 0x30, 0x39,
	0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x22, 0x8a, 0x01, 0x0a, 0x11, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x48, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x49, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x50, 0x72, 0x65, 0x53, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x4b, 0x65, 0x79, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4b,
	0x65, 0x79, 0x12, 0x2b, 0x0a, 0x04, 0x6a, 0x77, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x6a, 0x77, 0x6b, 0x73, 0x22,
	0x74, 0x0a, 0x1c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12,
	0x3e, 0x0a, 0x1b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x73,
	0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x19, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x84, 0x01, 0x0a, 0x12, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6e, 0x0a, 0x1e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x62, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x1c,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x22, 0x51, 0x0a, 0x17,
	0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x1a, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c,
	0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x67, 0x64, 0x2d, 0x64,
	0x65, 0x76, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescOnce sync.Once
	file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescData = file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDesc
)

func file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescGZIP() []byte {
	file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescOnce.Do(func() {
		file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescData)
	})
	return file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDescData
}

var file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_github_com_plgd_dev_client_application_pb_initialize_proto_goTypes = []interface{}{
	(*InitializePreSharedKey)(nil),       // 0: service.pb.InitializePreSharedKey
	(*InitializeX509)(nil),               // 1: service.pb.InitializeX509
	(*InitializeRequest)(nil),            // 2: service.pb.InitializeRequest
	(*IdentityCertificateChallenge)(nil), // 3: service.pb.IdentityCertificateChallenge
	(*InitializeResponse)(nil),           // 4: service.pb.InitializeResponse
	(*FinishInitializeRequest)(nil),      // 5: service.pb.FinishInitializeRequest
	(*FinishInitializeResponse)(nil),     // 6: service.pb.FinishInitializeResponse
	(*structpb.Struct)(nil),              // 7: google.protobuf.Struct
}
var file_github_com_plgd_dev_client_application_pb_initialize_proto_depIdxs = []int32{
	0, // 0: service.pb.InitializeRequest.pre_shared_key:type_name -> service.pb.InitializePreSharedKey
	7, // 1: service.pb.InitializeRequest.jwks:type_name -> google.protobuf.Struct
	3, // 2: service.pb.InitializeResponse.identity_certificate_challenge:type_name -> service.pb.IdentityCertificateChallenge
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_github_com_plgd_dev_client_application_pb_initialize_proto_init() }
func file_github_com_plgd_dev_client_application_pb_initialize_proto_init() {
	if File_github_com_plgd_dev_client_application_pb_initialize_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitializePreSharedKey); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitializeX509); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitializeRequest); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentityCertificateChallenge); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitializeResponse); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishInitializeRequest); i {
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
		file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishInitializeResponse); i {
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
			RawDescriptor: file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_plgd_dev_client_application_pb_initialize_proto_goTypes,
		DependencyIndexes: file_github_com_plgd_dev_client_application_pb_initialize_proto_depIdxs,
		MessageInfos:      file_github_com_plgd_dev_client_application_pb_initialize_proto_msgTypes,
	}.Build()
	File_github_com_plgd_dev_client_application_pb_initialize_proto = out.File
	file_github_com_plgd_dev_client_application_pb_initialize_proto_rawDesc = nil
	file_github_com_plgd_dev_client_application_pb_initialize_proto_goTypes = nil
	file_github_com_plgd_dev_client_application_pb_initialize_proto_depIdxs = nil
}
