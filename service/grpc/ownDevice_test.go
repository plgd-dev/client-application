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
	"os"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerOwnDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewClientApplicationGetDevicesServer(ctx))
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

func TestClientApplicationServerOwnDeviceViaManufacturerCertificate(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	cfg := test.MakeDeviceConfig()
	cfg.COAP.OwnershipTransfer.Methods = []device.OwnershipTransferMethod{device.OwnershipTransferManufacturerCertificate}
	cfg.COAP.OwnershipTransfer.Manufacturer.TLS.CAPool = os.Getenv("MFG_ROOT_CA_CRT")
	cfg.COAP.OwnershipTransfer.Manufacturer.TLS.CertFile = os.Getenv("MFG_CLIENT_APPLICATION_CRT")
	cfg.COAP.OwnershipTransfer.Manufacturer.TLS.KeyFile = os.Getenv("MFG_CLIENT_APPLICATION_KEY")
	s, teardown, err := test.NewClientApplicationServer(ctx, test.WithDeviceConfig(cfg))
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewClientApplicationGetDevicesServer(ctx))
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
