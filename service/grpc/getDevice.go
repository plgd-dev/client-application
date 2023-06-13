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
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) GetDevice(ctx context.Context, req *pb.GetDeviceRequest) (*grpcgwPb.Device, error) {
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
	devLinks := links.GetResourceLinks(plgdDevice.ResourceType)
	if len(devLinks) == 0 {
		return nil, status.Errorf(codes.NotFound, "cannot find device resource %v at device %v", plgdDevice.ResourceType, req.GetDeviceId())
	}
	_, err = getResourceAndRefreshCache(ctx, dev, devLinks[0], "")
	if err != nil {
		return nil, err
	}

	return dev.ToProto(), nil
}
