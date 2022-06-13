package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeviceGatewayServerUpdateResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	type args struct {
		req *grpcgwPb.UpdateResourceRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *grpcgwPb.UpdateResourceResponse
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "doxm update",
			args: args{
				req: &grpcgwPb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     doxm.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"oxmsel":0}`),
					},
				},
			},
			want: &grpcgwPb.UpdateResourceResponse{
				Data: &events.ResourceUpdated{
					Content: &commands.Content{},
					Status:  commands.Status_OK,
				},
			},
		},
		{
			name: "device resource - fail",
			args: args{
				req: &grpcgwPb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     device.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.PermissionDenied,
		},
		{
			name: "unknown device",
			args: args{
				req: &grpcgwPb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: uuid.NewString(),
						Href:     device.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "unknown href",
			args: args{
				req: &grpcgwPb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     "/unknown",
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "permission denied - cannot establish TLS connection",
			args: args{
				req: &grpcgwPb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     "/light/1",
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.PermissionDenied,
		},
	}

	s, teardown, err := test.NewDeviceGatewayServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewDeviceGatewayGetDevicesServer(ctx))
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateResource(ctx, tt.args.req)
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
