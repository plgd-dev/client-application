package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) DisownDevice(ctx context.Context, req *pb.DisownDeviceRequest) (*pb.DisownDeviceResponse, error) {
	dev, err := s.getDevice(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	if dev.ToProto().OwnershipStatus != grpcgwPb.Device_OWNED {
		return nil, status.Error(codes.PermissionDenied, "device is not owned")
	}

	err = dev.Disown(ctx, links)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot disown device %v: %w", dev.ID, err)).Err()
	}
	err = s.deleteDevice(ctx, dev.ID)
	if err != nil {
		log.Errorf("cannot remove device %v from cache: %v", dev.ID, err)
	}

	return &pb.DisownDeviceResponse{}, nil
}
