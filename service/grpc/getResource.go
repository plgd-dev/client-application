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
	"google.golang.org/grpc/status"
)

func getResourceAndRefreshCache(ctx context.Context, dev *device, link schema.ResourceLink) (*commands.Content, error) {
	codec := rawcodec.GetRawCodec(message.AppOcfCbor)
	var data []byte

	if dev.ToProto().OwnershipStatus != grpcgwPb.Device_OWNED && len(link.Endpoints.FilterUnsecureEndpoints()) == 0 {
		return nil, status.Error(codes.PermissionDenied, "device is not owned")
	}
	err := dev.GetResourceWithCodec(ctx, link, codec, &data)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", link.Href, dev.ID, err)).Err()
	}
	contentType := ""
	if len(data) > 0 {
		contentType = message.AppOcfCbor.String()
	}
	content := &commands.Content{
		ContentType: contentType,
		Data:        data,
	}
	// we update device resource body only for device resource
	if strings.Contains(link.ResourceTypes, plgdDevice.ResourceType) {
		dev.updateDeviceResourceBody(content)
	}
	return content, nil
}

func (s *DeviceGatewayServer) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*grpcgwPb.Resource, error) {
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLink(ctx, req.GetResourceId())
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