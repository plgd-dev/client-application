package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	hubTest "github.com/plgd-dev/hub/v2/test"
	hubTestPb "github.com/plgd-dev/hub/v2/test/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeviceGatewayServerGetDeviceResourceLinks(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	type args struct {
		req *pb.GetDeviceResourceLinksRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *events.ResourceLinksPublished
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "device links",
			args: args{
				req: &pb.GetDeviceResourceLinksRequest{
					DeviceId: dev.Id,
				},
			},
			want: hubTestPb.CleanUpResourceLinksPublished(&events.ResourceLinksPublished{
				DeviceId:  dev.Id,
				Resources: commands.SchemaResourceLinksToResources(test.GetDeviceResourceLinks(), time.Time{}),
			}, true),
		},
		{
			name: "unknown device",
			args: args{
				req: &pb.GetDeviceResourceLinksRequest{
					DeviceId: uuid.NewString(),
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
	}

	s, err := test.NewDeviceGatewayServer(ctx)
	require.NoError(t, err)
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewDeviceGatewayGetDevicesServer(ctx))
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetDeviceResourceLinks(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrCode.String(), status.Code(err).String())
				return
			}
			require.NoError(t, err)
			got = hubTestPb.CleanUpResourceLinksPublished(got, true)
			for _, r := range got.GetResources() {
				require.NotEmpty(t, r.GetEndpointInformations())
				r.EndpointInformations = nil
				require.NotEmpty(t, r.GetPolicy())
				r.Policy = nil
				require.NotEmpty(t, r.GetAnchor())
				r.Anchor = ""
			}
			require.NoError(t, err)
			hubTest.CheckProtobufs(t, tt.want, got, hubTest.RequireToCheckFunc(require.Equal))
		})
	}
}
