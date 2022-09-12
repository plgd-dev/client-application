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
	"github.com/plgd-dev/client-application/pkg/rawcodec"
	"github.com/plgd-dev/go-coap/v2/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
)

func (s *ClientApplicationServer) DeleteResource(ctx context.Context, req *pb.DeleteResourceRequest) (*grpcgwPb.DeleteResourceResponse, error) {
	devID, err := strDeviceID2UUID(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(devID)
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLinkAndCheckAccess(ctx, req.GetResourceId())
	if err != nil {
		return nil, err
	}
	var response []byte
	err = dev.DeleteResourceWithCodec(ctx, link, rawcodec.GetRawCodec(message.AppOcfCbor), &response)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot delete resource %v for device %v: %w", link.Href, dev.ID, err)).Err()
	}
	return &grpcgwPb.DeleteResourceResponse{
		Data: &events.ResourceDeleted{
			Content: responseToData(response),
			Status:  commands.Status_OK,
		},
	}, nil
}
