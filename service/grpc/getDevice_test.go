// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/v2/schema/configuration"
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{
		UseMulticast: []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV6},
	}, test.NewClientApplicationGetDevicesServer(ctx))
	require.NoError(t, err)

	time.Sleep(time.Second)

	d1, err := s.GetDevice(ctx, &pb.GetDeviceRequest{
		DeviceId: dev.GetId(),
	})
	require.NoError(t, err)
	require.Equal(t, dev.GetData(), d1.GetData())

	_, err = s.OwnDevice(ctx, &pb.OwnDeviceRequest{
		DeviceId: dev.GetId(),
	})
	require.NoError(t, err)

	newName := test.DevsimName + "_new"
	_, err = s.UpdateResource(ctx, &pb.UpdateResourceRequest{
		ResourceId: commands.NewResourceID(dev.GetId(), configuration.ResourceURI),
		Content: &grpcgwPb.Content{
			ContentType: serviceHttp.ApplicationJsonContentType,
			Data:        []byte(`{"n":"` + newName + `"}`),
		},
	})
	require.NoError(t, err)

	d1, err = s.GetDevice(ctx, &pb.GetDeviceRequest{
		DeviceId: dev.GetId(),
	})
	require.NoError(t, err)
	var v plgdDevice.Device
	err = cbor.Decode(d1.GetData().GetContent().GetData(), &v)
	require.NoError(t, err)
	require.Equal(t, newName, v.Name)

	_, err = s.UpdateResource(ctx, &pb.UpdateResourceRequest{
		ResourceId: commands.NewResourceID(dev.GetId(), configuration.ResourceURI),
		Content: &grpcgwPb.Content{
			ContentType: serviceHttp.ApplicationJsonContentType,
			Data:        []byte(`{"n":"` + test.DevsimName + `"}`),
		},
	})
	require.NoError(t, err)

	_, err = s.DisownDevice(ctx, &pb.DisownDeviceRequest{
		DeviceId: dev.GetId(),
	})
	require.NoError(t, err)
}
