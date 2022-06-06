package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	serviceGrpc "github.com/plgd-dev/client-application/service/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type testDeviceGatewayGetDevicesServer struct {
	grpc.ServerStream
	devices []*pb.Device
}

func (s *testDeviceGatewayGetDevicesServer) Send(d *pb.Device) error {
	s.devices = append(s.devices, d)
	return nil
}

func (s *testDeviceGatewayGetDevicesServer) Context() context.Context {
	return context.Background()
}

func TestDeviceGatewayServer_GetDevices(t *testing.T) {
	type args struct {
		req *pb.GetDevicesRequest
		srv pb.DeviceGateway_GetDevicesServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				req: &pb.GetDevicesRequest{},
				srv: &testDeviceGatewayGetDevicesServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serviceGrpc.DeviceGatewayServer{}
			err := s.GetDevices(tt.args.req, tt.args.srv)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, d := range tt.args.srv.(*testDeviceGatewayGetDevicesServer).devices {
				fmt.Printf("%v\n", d)
			}
		})
	}
}
