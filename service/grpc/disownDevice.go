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
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) DisownDevice(ctx context.Context, req *pb.DisownDeviceRequest) (*pb.DisownDeviceResponse, error) {
	dev, err := s.getDevice(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	if dev.ToProto().OwnershipStatus != grpcgwPb.Device_OWNED {
		return nil, status.Error(codes.PermissionDenied, "device is not owned")
	}

	err = dev.Disown(ctx, links)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot disown device %v: %w", dev.ID, err)).Err()
	}
	err = s.deleteDevice(ctx, dev.ID)
	if err != nil {
		log.Errorf("cannot remove device %v from cache: %v", dev.ID, err)
	}

	return &pb.DisownDeviceResponse{}, nil
}
