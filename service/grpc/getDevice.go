package grpc

import (
	"context"

	"github.com/plgd-dev/client-application/pb"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DeviceGatewayServer) GetDevice(ctx context.Context, req *pb.GetDeviceRequest) (*grpcgwPb.Device, error) {
	dev, err := s.getDevice(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	devLinks := links.GetResourceLinks(plgdDevice.ResourceType)
	if len(devLinks) == 0 {
		return nil, status.Errorf(codes.NotFound, "cannot find device resource %v at device %v", plgdDevice.ResourceType, req.GetDeviceId())
	}
	_, err = getResourceAndRefreshCache(ctx, dev, devLinks[0])
	if err != nil {
		return nil, err
	}

	return dev.ToProto(), nil
}
