package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/rawcodec"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/device/schema/interfaces"
	"github.com/plgd-dev/go-coap/v2/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
)

func (s *ClientApplicationServer) CreateResource(ctx context.Context, req *pb.CreateResourceRequest) (*grpcgwPb.CreateResourceResponse, error) {
	createData, err := convContentToOcfCbor(req.GetContent())
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
	err = dev.UpdateResourceWithCodec(ctx, link, rawcodec.GetRawCodec(message.AppOcfCbor), createData, &response, coap.WithInterface(interfaces.OC_IF_CREATE))
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot create resource %v for device %v: %w", link.Href, dev.ID, err)).Err()
	}
	return &grpcgwPb.CreateResourceResponse{
		Data: &events.ResourceCreated{
			Content: responseToData(response),
			Status:  commands.Status_OK,
		},
	}, nil
}
