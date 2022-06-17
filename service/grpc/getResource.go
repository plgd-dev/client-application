package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/rawcodec"
	"github.com/plgd-dev/device/schema"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/go-coap/v2/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/strings"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
)

func getResourceAndRefreshCache(ctx context.Context, dev *device, link schema.ResourceLink) (*commands.Content, error) {
	var response []byte
	err := dev.GetResourceWithCodec(ctx, link, rawcodec.GetRawCodec(message.AppOcfCbor), &response)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", link.Href, dev.ID, err)).Err()
	}
	content := responseToData(response)
	// we update device resource body only for device resource
	if strings.Contains(link.ResourceTypes, plgdDevice.ResourceType) {
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
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLinkAndCheckAccess(ctx, req.GetResourceId())
	if err != nil {
		return nil, err
	}
	content, err := getResourceAndRefreshCache(ctx, dev, link)
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
