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
	"slices"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/rawcodec"
	pkgCoap "github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/device/v2/schema"
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	"github.com/plgd-dev/go-coap/v3/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
)

func getResourceAndRefreshCache(ctx context.Context, dev *device, link schema.ResourceLink, resourceInterface string) (*commands.Content, error) {
	var response []byte
	options := make([]func(message.Options) message.Options, 0, 2)
	options = append(options, pkgCoap.WithDeviceID(dev.DeviceID()))
	if resourceInterface != "" {
		options = append(options, pkgCoap.WithInterface(resourceInterface))
	}

	err := dev.GetResourceWithCodec(ctx, link, rawcodec.GetRawCodec(message.AppOcfCbor), &response, options...)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", link.Href, dev.ID, err)).Err()
	}
	content := responseToData(response)
	// we update device resource body only for device resource
	if slices.Contains(link.ResourceTypes, plgdDevice.ResourceType) && resourceInterface == "" {
		dev.updateDeviceResourceBody(content)
	}
	return content, nil
}

func responseToData(response []byte) *commands.Content {
	contentType := ""
	if len(response) > 0 {
		contentType = message.AppOcfCbor.String()
	}
	return &commands.Content{
		ContentType: contentType,
		Data:        response,
	}
}

func (s *ClientApplicationServer) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*grpcgwPb.Resource, error) {
	devID, err := strDeviceID2UUID(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(devID)
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLinkAndCheckAccess(ctx, req.GetResourceId(), req.GetResourceInterface())
	if err != nil {
		return nil, err
	}
	content, err := getResourceAndRefreshCache(ctx, dev, link, req.GetResourceInterface())
	if err != nil {
		return nil, err
	}
	return &grpcgwPb.Resource{
		Data: &events.ResourceChanged{
			Content: content,
			Status:  commands.Status_OK,
		},
		Types: link.ResourceTypes,
	}, nil
}
