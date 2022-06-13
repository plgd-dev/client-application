package grpc

import (
	"context"
	"sync"

	"github.com/plgd-dev/client-application/pb"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceGatewayServer struct {
	serviceDevice *serviceDevice.Service
	logger        log.Logger
	devices       sync.Map
	pb.UnimplementedDeviceGatewayServer
}

func NewDeviceGatewayServer(serviceDevice *serviceDevice.Service, logger log.Logger) *DeviceGatewayServer {
	return &DeviceGatewayServer{
		serviceDevice: serviceDevice,
		logger:        logger,
	}
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

func (s *DeviceGatewayServer) deleteDevice(ctx context.Context, deviceID string) error {
	d, ok := s.devices.LoadAndDelete(deviceID)
	if !ok {
		return nil
	}
	dev, ok := d.(*device)
	if !ok {
		return status.Error(codes.Internal, "cast error")
	}
	if err := dev.Close(ctx); err != nil {
		return status.Errorf(codes.Internal, "cannot close device %v connections: %v", deviceID, err)
	}
	return nil
}
