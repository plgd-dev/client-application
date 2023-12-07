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
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/kit/v2/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) validateState(state uuid.UUID) (*serviceDevice.Service, bool) {
	item := s.csrCache.Get(state)
	if item == nil {
		return nil, false
	}
	s.csrCache.Delete(state)
	if item.IsExpired() {
		return nil, false
	}
	return item.Value(), true
}

func (s *ClientApplicationServer) signIdentityCertificateRemotely() bool {
	devService := s.serviceDevice.Load()
	if devService == nil {
		return false
	}
	return devService.GetDeviceAuthenticationMode() == pb.GetConfigurationResponse_X509
}

func (s *ClientApplicationServer) updateIdentityCertificate(ctx context.Context, req *pb.FinishInitializeRequest) error {
	state, err := uuid.Parse(req.GetState())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot parse state: %v", err)
	}
	devState, ok := s.validateState(state)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid state")
	}
	owner, err := grpc.OwnerFromTokenMD(ctx, s.GetConfig().RemoteProvisioning.GetJwtOwnerClaim())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "cannot get owner from token: %v", err)
	}
	ownerID := events.OwnerToUUID(owner)
	certs, err := security.ParseX509FromPEM(req.GetCertificate())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot parse certificate: %v", err)
	}
	ident, err := coap.GetDeviceIDFromIdentityCertificate(certs[0])
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot get owner id from certificate: %v", err)
	}
	if ownerID != ident {
		return status.Errorf(codes.InvalidArgument, "invalid owner id")
	}
	if err := devState.SetIdentityCertificate(owner, req.GetCertificate()); err != nil {
		return status.Errorf(codes.Internal, "cannot set certificate: %v", err)
	}
	s.init(ctx, devState)
	return nil
}
