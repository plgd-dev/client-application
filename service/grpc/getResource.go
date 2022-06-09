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

func (s *DeviceGatewayServer) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*grpcgwPb.Resource, error) {
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLink(ctx, req.GetResourceId())
	if err != nil {
		return nil, err
	}
	codec := rawcodec.GetRawCodec(message.AppOcfCbor)
	var data []byte
	err = dev.GetResourceWithCodec(ctx, link, codec, &data)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", req.GetResourceId().GetHref(), dev.ID, err)).Err()
	}
	contentType := ""
	if len(data) > 0 {
		contentType = message.AppOcfCbor.String()
	}
	return &grpcgwPb.Resource{
		Data: &events.ResourceChanged{
			Content: &commands.Content{
				ContentType: contentType,
				Data:        data,
			},
			Status: commands.Status_OK,
		},
		Types: link.ResourceTypes,
	}, nil
}
