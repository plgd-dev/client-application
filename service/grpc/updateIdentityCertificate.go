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
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/kit/v2/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) validateState(state uuid.UUID) bool {
	item := s.csrCache.Get(state)
	if item == nil {
		return false
	}
	s.csrCache.Delete(state)
	return !item.IsExpired()
}

func (s *ClientApplicationServer) UpdateIdentityCertificate(ctx context.Context, req *pb.UpdateIdentityCertificateRequest) (*pb.UpdateIdentityCertificateResponse, error) {
	if s.remoteProvisioningConfig.Mode != remoteProvisioning.Mode_UserAgent {
		return nil, status.Errorf(codes.Unimplemented, "not supported")
	}
	state, err := uuid.Parse(req.State)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse state: %v", err)
	}
	if !s.validateState(state) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid state")
	}
	owner, err := grpc.OwnerFromTokenMD(ctx, s.remoteProvisioningConfig.Authorization.OwnerClaim)
	if err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.Unauthenticated, "cannot get owner from token: %v", err))
	}
	owner = events.OwnerToUUID(owner)
	certs, err := security.ParseX509FromPEM([]byte(req.Certificate))
	if err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.InvalidArgument, "cannot parse certificate: %v", err))
	}
	ident, err := coap.GetDeviceIDFromIdentityCertificate(certs[0])
	if err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.InvalidArgument, "cannot get owner id from certificate: %v", err))
	}
	if owner != ident {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.InvalidArgument, "invalid owner id"))
	}
	if err := s.serviceDevice.SetIdentityCertificate([]byte(req.Certificate)); err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.Internal, "cannot set certificate: %v", err))
	}
	return &pb.UpdateIdentityCertificateResponse{}, nil
}
