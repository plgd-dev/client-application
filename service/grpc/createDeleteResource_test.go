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
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/go-coap/v2/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerCreateDeleteResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{
		UseMulticast: []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4},
	}, test.NewClientApplicationGetDevicesServer(ctx))
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

	_, err = s.CreateResource(ctx, &pb.CreateResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, hubTest.TestResourceSwitchesHref),
		Content: &grpcgwPb.Content{
			ContentType: message.AppOcfCbor.String(),
			Data:        hubTest.EncodeToCbor(t, hubTest.MakeSwitchResourceDefaultData()),
		},
	})
	require.NoError(t, err)

	_, err = s.DeleteResource(ctx, &pb.DeleteResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, hubTest.TestResourceSwitchesInstanceHref("1")),
	})
	require.NoError(t, err)

	// duplicity delete
	_, err = s.DeleteResource(ctx, &pb.DeleteResourceRequest{
		ResourceId: commands.NewResourceID(dev.Id, hubTest.TestResourceSwitchesInstanceHref("1")),
	})
	require.Error(t, err)

	_, err = s.DisownDevice(ctx, &pb.DisownDeviceRequest{
		DeviceId: dev.Id,
	})
	require.NoError(t, err)
}
