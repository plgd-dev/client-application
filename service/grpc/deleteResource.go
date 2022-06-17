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
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
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
