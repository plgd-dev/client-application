package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/configuration"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	s, teardown, err := test.NewDeviceGatewayServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{
		UseMulticast: []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4},
	}, test.NewDeviceGatewayGetDevicesServer(ctx))
	require.NoError(t, err)

	d1, err := s.GetDevice(ctx, &pb.GetDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)
	require.Equal(t, dev, d1)

	_, err = s.OwnDevice(ctx, &pb.OwnDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)

	newName := test.DevsimName + "_new"
	_, err = s.UpdateResource(ctx, &grpcgwPb.UpdateResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, configuration.ResourceURI),
		Content: &grpcgwPb.Content{
			ContentType: serviceHttp.ApplicationJsonContentType,
			Data:        []byte(`{"n":"` + newName + `"}`),
		},
	})
	require.NoError(t, err)

	d1, err = s.GetDevice(ctx, &pb.GetDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)
	var v plgdDevice.Device
	err = cbor.Decode(d1.GetData().GetContent().GetData(), &v)
	require.NoError(t, err)
	require.Equal(t, newName, v.Name)

	_, err = s.UpdateResource(ctx, &grpcgwPb.UpdateResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, configuration.ResourceURI),
		Content: &grpcgwPb.Content{
			ContentType: serviceHttp.ApplicationJsonContentType,
			Data:        []byte(`{"n":"` + test.DevsimName + `"}`),
		},
	})
	require.NoError(t, err)

	_, err = s.DisownDevice(ctx, &pb.DisownDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)
}
