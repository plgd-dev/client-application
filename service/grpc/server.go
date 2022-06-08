package grpc

import (
	"sync"

	"github.com/plgd-dev/client-application/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceGatewayServer struct {
	devices sync.Map
	pb.UnimplementedDeviceGatewayServer
}

func NewDeviceGatewayServer() *DeviceGatewayServer {
	return &DeviceGatewayServer{}
}

func (s *DeviceGatewayServer) getDevice(deviceID string) (*device, error) {
	d, ok := s.devices.Load(deviceID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "device %v not found", deviceID)
	}
	dev, ok := d.(*device)
	if !ok {
		return nil, status.Error(codes.Internal, "cast error")
	}
	return dev, nil
}
