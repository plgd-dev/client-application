package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceGrpc "github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/device"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeviceGatewayServerGetResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	type args struct {
		req *pb.GetResourceRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *pb.GetResourceResponse
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "device resource",
			args: args{
				req: &pb.GetResourceRequest{
					ResourceId: &pb.ResourceId{
						DeviceId: dev.Id,
						Href:     device.ResourceURI,
					},
				},
			},
			want: &pb.GetResourceResponse{
				Content: dev.Content,
			},
		},
		{
			name: "unknown device",
			args: args{
				req: &pb.GetResourceRequest{
					ResourceId: &pb.ResourceId{
						DeviceId: uuid.NewString(),
						Href:     device.ResourceURI,
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "unknown href",
			args: args{
				req: &pb.GetResourceRequest{
					ResourceId: &pb.ResourceId{
						DeviceId: dev.Id,
						Href:     "/unknown",
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "unavailable - cannot establish TLS connection",
			args: args{
				req: &pb.GetResourceRequest{
					ResourceId: &pb.ResourceId{
						DeviceId: dev.Id,
						Href:     "/light/1",
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.Unavailable,
		},
	}

	s := &serviceGrpc.DeviceGatewayServer{}
	err := s.GetDevices(&pb.GetDevicesRequest{}, test.NewDeviceGatewayGetDevicesServer(ctx))
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetResource(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrCode.String(), status.Code(err).String())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
