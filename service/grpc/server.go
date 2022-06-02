package grpc

import (
	"github.com/plgd-dev/client-application/pb"
)

type DeviceGatewayServer struct {
	pb.UnimplementedDeviceGatewayServer
}

func NewDeviceGatewayServer() *DeviceGatewayServer {
	return &DeviceGatewayServer{}
}
