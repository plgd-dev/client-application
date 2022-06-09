package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"google.golang.org/grpc/codes"
)

func (s *DeviceGatewayServer) OwnDevice(ctx context.Context, req *pb.OwnDeviceRequest) (*pb.OwnDeviceResponse, error) {
	dev, err := s.getDevice(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinks(ctx)
	if err != nil {
		return nil, err
	}

	err = dev.Own(ctx, links, s.serviceDevice.GetJustWorksClient(), s.serviceDevice.GetOwnOptions()...)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot own device %v: %w", dev.ID, err)).Err()
	}
	dev.updateOwnershipStatus(grpcgwPb.Device_OWNED)

	return &pb.OwnDeviceResponse{}, nil
}
