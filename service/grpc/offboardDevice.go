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
	"github.com/plgd-dev/device/v2/schema/cloud"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) OffboardDevice(ctx context.Context, req *pb.OffboardDeviceRequest) (resp *pb.OffboardDeviceResponse, err error) {
	devID, err := strDeviceID2UUID(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(devID)
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	cloudLinks := links.GetResourceLinks(cloud.ResourceType)
	if len(cloudLinks) == 0 {
		return nil, status.Errorf(codes.NotFound, "cannot find cloud resource for device %v", devID)
	}
	cloudLink := cloudLinks[0]
	if err = dev.checkAccess(cloudLink); err != nil {
		return nil, err
	}
	err = dev.UpdateResource(ctx, cloudLink, cloud.ConfigurationUpdateRequest{}, nil)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot update resource %v for device %v: %w", cloudLink.Href, dev.ID, err)).Err()
	}
	return &pb.OffboardDeviceResponse{}, nil
}
