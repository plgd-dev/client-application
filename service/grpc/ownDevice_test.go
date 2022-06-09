package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerOwnDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	s, err := test.NewDeviceGatewayServer(ctx)
	require.NoError(t, err)
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewDeviceGatewayGetDevicesServer(ctx))
	require.NoError(t, err)

	_, err = s.OwnDevice(ctx, &pb.OwnDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)

	_, err = s.GetResource(ctx, &pb.GetResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, "/light/1"),
	})
	require.NoError(t, err)

	_, err = s.DisownDevice(ctx, &pb.DisownDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)
}
