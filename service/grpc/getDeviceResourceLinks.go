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
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
)

func (s *ClientApplicationServer) GetDeviceResourceLinks(ctx context.Context, req *pb.GetDeviceResourceLinksRequest) (*events.ResourceLinksPublished, error) {
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
	return &events.ResourceLinksPublished{
		Resources: commands.SchemaResourceLinksToResources(links, time.Time{}),
		DeviceId:  dev.DeviceID(),
	}, nil
}
