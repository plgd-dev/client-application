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

	"github.com/plgd-dev/client-application/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) Initialize(ctx context.Context, req *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	if s.serviceDevice.IsInitialized() {
		return nil, status.Errorf(codes.AlreadyExists, "already initialized")
	}
	if !s.signIdentityCertificateRemotely() {
		return nil, status.Errorf(codes.InvalidArgument, "initialize with certificate is disabled")
	}
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

func (s *ClientApplicationServer) FinishInitialize(ctx context.Context, req *pb.FinishInitializeRequest) (*pb.FinishInitializeResponse, error) {
	if s.serviceDevice.IsInitialized() {
		return nil, status.Errorf(codes.AlreadyExists, "already initialized")
	}
	if err := s.updateIdentityCertificate(ctx, req); err != nil {
		return nil, err
	}
	return &pb.FinishInitializeResponse{}, nil
}
