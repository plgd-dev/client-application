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

package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errAlreadyInitialized = "already initialized"

func (s *ClientApplicationServer) InitializeRemoteProvisioning(ctx context.Context, req *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	err := s.UpdateJSONWebKeys(ctx, req.GetJwks())
	if err != nil {
		return nil, err
	}
	respCsr, err := s.getIdentityCSR(ctx)
	if err != nil {
		s.reset(ctx)
		return nil, err
	}
	return &pb.InitializeResponse{
		IdentityCertificateChallenge: respCsr,
	}, nil
}

func (s *ClientApplicationServer) UpdatePSK(subjectUUID, key string) error {
	s.updatePSKLock.Lock()
	defer s.updatePSKLock.Unlock()
	isInitialized := s.serviceDevice.IsInitialized()
	if isInitialized && (subjectUUID != "" || key != "") {
		return status.Errorf(codes.FailedPrecondition, errAlreadyInitialized)
	}
	if !isInitialized && (subjectUUID == "" && key == "") {
		return status.Errorf(codes.FailedPrecondition, "not initialized")
	}
	cfg := s.GetConfig()
	cfg.Clients.Device.COAP.TLS.PreSharedKey.Key = key
	cfg.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr = subjectUUID
	return s.StoreConfig(cfg)
}

func (s *ClientApplicationServer) Initialize(ctx context.Context, req *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	if s.serviceDevice.IsInitialized() {
		return nil, status.Errorf(codes.FailedPrecondition, errAlreadyInitialized)
	}
	if s.signIdentityCertificateRemotely() {
		return s.InitializeRemoteProvisioning(ctx, req)
	}
	if req.GetPreSharedKey().GetSubjectId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared subjectUuid(%v)", req.GetPreSharedKey().GetSubjectId())
	}
	_, err := uuid.Parse(req.GetPreSharedKey().GetSubjectId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared subjectUuid(%v): %v", req.GetPreSharedKey().GetSubjectId(), err)
	}
	if req.GetPreSharedKey().GetKey() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared key")
	}
	if err := s.UpdatePSK(req.GetPreSharedKey().GetSubjectId(), req.GetPreSharedKey().GetKey()); err != nil {
		return nil, err
	}
	if _, err := s.ClearCache(ctx, &pb.ClearCacheRequest{}); err != nil {
		log.Warnf("cannot clear device cache: %v", err)
	}
	return &pb.InitializeResponse{}, nil
}

func (s *ClientApplicationServer) FinishInitialize(ctx context.Context, req *pb.FinishInitializeRequest) (*pb.FinishInitializeResponse, error) {
	if s.serviceDevice.IsInitialized() {
		return nil, status.Errorf(codes.FailedPrecondition, errAlreadyInitialized)
	}
	if err := s.updateIdentityCertificate(ctx, req); err != nil {
		return nil, err
	}
	if _, err := s.ClearCache(ctx, &pb.ClearCacheRequest{}); err != nil {
		log.Warnf("cannot clear device cache: %v", err)
	}
	return &pb.FinishInitializeResponse{}, nil
}
