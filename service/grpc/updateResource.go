//***************************************************************************
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
//*************************************************************************

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
	"github.com/plgd-dev/kit/v2/codec/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convContentToOcfCbor(content *grpcgwPb.Content) ([]byte, error) {
	switch content.GetContentType() {
	case message.AppCBOR.String(), message.AppOcfCbor.String():
		return content.GetData(), nil
	case message.AppJSON.String():
		data, err := json.ToCBOR(string(content.GetData()))
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot convert json to cbor: %v", err)
		}
		return data, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "unsupported content type '%v'", content.GetContentType())
}

func (s *ClientApplicationServer) UpdateResource(ctx context.Context, req *pb.UpdateResourceRequest) (*grpcgwPb.UpdateResourceResponse, error) {
	updateData, err := convContentToOcfCbor(req.GetContent())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLinkAndCheckAccess(ctx, req.GetResourceId())
	if err != nil {
		return nil, err
	}
	var response []byte
	err = dev.UpdateResourceWithCodec(ctx, link, rawcodec.GetRawCodec(message.AppOcfCbor), updateData, &response)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", req.GetResourceId().GetHref(), dev.ID, err)).Err()
	}
	return &grpcgwPb.UpdateResourceResponse{
		Data: &events.ResourceUpdated{
			Content: responseToData(response),
			Status:  commands.Status_OK,
		},
	}, nil
}
