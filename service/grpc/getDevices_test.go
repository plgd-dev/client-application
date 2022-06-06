package grpc_test

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceGrpc "github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetDevices(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	u, err := url.Parse(device.Endpoints[0])
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	type args struct {
		req *pb.GetDevicesRequest
		srv pb.DeviceGateway_GetDevicesServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*pb.Device
	}{
		{
			name: "by multicast",
			args: args{
				req: &pb.GetDevicesRequest{
					UseMulticast: []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4},
				},
				srv: test.NewDeviceGatewayGetDevicesServer(ctx),
			},
			want: []*pb.Device{
				device,
			},
		},
		{
			name: "by ip",
			args: args{
				req: &pb.GetDevicesRequest{
					UseEndpoints: []string{u.Hostname()},
				},
				srv: test.NewDeviceGatewayGetDevicesServer(ctx),
			},
			want: []*pb.Device{
				device,
			},
		},
		{
			name: "by ip:port",
			args: args{
				req: &pb.GetDevicesRequest{
					UseEndpoints: []string{u.Host},
				},
				srv: test.NewDeviceGatewayGetDevicesServer(ctx),
			},
			want: []*pb.Device{
				device,
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
			got := tt.args.srv.(*test.DeviceGatewayGetDevicesServer).Devices
			require.NotEmpty(t, got)
			require.Len(t, got[0].Endpoints, 4)
			require.True(t, strings.Contains(got[0].Endpoints[0], "coap://"))
			require.True(t, strings.Contains(got[0].Endpoints[1], "coap+tcp://"))
			require.True(t, strings.Contains(got[0].Endpoints[2], "coaps://"))
			require.True(t, strings.Contains(got[0].Endpoints[3], "coaps+tcp://"))
			assert.Equal(t, tt.want, got)
		})
	}
}
