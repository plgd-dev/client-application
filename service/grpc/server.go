package grpc

import (
	"sync"

	"github.com/plgd-dev/client-application/pb"
)

type DeviceGatewayServer struct {
	devices sync.Map
	pb.UnimplementedDeviceGatewayServer
}

func NewDeviceGatewayServer() *DeviceGatewayServer {
	return &DeviceGatewayServer{}
}
