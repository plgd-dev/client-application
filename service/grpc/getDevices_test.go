package grpc_test

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetDevices(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	u, err := url.Parse(device.Endpoints[0])
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	device.OwnershipStatus = grpcgwPb.Device_UNOWNED

	type args struct {
		req *pb.GetDevicesRequest
		srv pb.DeviceGateway_GetDevicesServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*grpcgwPb.Device
	}{
		{
			name: "by multicast",
			args: args{
				req: &pb.GetDevicesRequest{
					UseMulticast: []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4},
				},
				srv: test.NewDeviceGatewayGetDevicesServer(ctx),
			},
			want: []*grpcgwPb.Device{
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
			want: []*grpcgwPb.Device{
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
			want: []*grpcgwPb.Device{
				device,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, teardown, err := test.NewDeviceGatewayServer(ctx)
			require.NoError(t, err)
			defer teardown()
			err = s.GetDevices(tt.args.req, tt.args.srv)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			server, ok := tt.args.srv.(*test.DeviceGatewayGetDevicesServer)
			require.True(t, ok)
			got := server.Devices
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
