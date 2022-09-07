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
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) GetIdentityCSR(ctx context.Context, req *pb.GetIdentityCSRRequest) (*pb.GetIdentityCSRResponse, error) {
	if s.remoteProvisioningConfig.Mode != remoteProvisioning.Mode_UserAgent {
		return nil, status.Errorf(codes.Unimplemented, "not supported")
	}
	owner, err := grpc.OwnerFromTokenMD(ctx, s.remoteProvisioningConfig.Authorization.OwnerClaim)
	if err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.Unauthenticated, "cannot get owner from token: %v", err))
	}
	csr, err := s.serviceDevice.GetIdentityCSR(events.OwnerToUUID(owner))
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	state := uuid.New()
	s.csrCache.Set(state, true, s.remoteProvisioningConfig.UserAgentConfig.CSRChallengeStateExpiration)
	return &pb.GetIdentityCSRResponse{
		CertificateSigningRequest: string(csr),
		State:                     state.String(),
	}, nil
}
