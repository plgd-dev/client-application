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
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerCheckForClosingInactivityConnection(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()
	run := func() {
		srv := test.NewClientApplicationGetDevicesServer(ctx)
		err = s.GetDevices(&pb.GetDevicesRequest{}, srv)
		require.NoError(t, err)
		require.Greater(t, len(srv.Devices), 0)
		time.Sleep(time.Second)
		numParallel := 10
		var wg sync.WaitGroup
		wg.Add(numParallel)
		fmt.Printf("%v test starts\n", time.Now())
		errCh := make(chan error, numParallel)
		for i := 0; i < numParallel; i++ {
			go func(i int) {
				defer wg.Done()
				fmt.Printf("%v GetDevice %v start\n", time.Now(), i)
				_, err := s.GetDevice(ctx, &pb.GetDeviceRequest{DeviceId: device.Id})
				fmt.Printf("%v GetDevice %v end %v\n", time.Now(), i, err)
				errCh <- err
			}(i)
		}
		wg.Wait()
		fmt.Printf("%v test ends\n", time.Now())
		close(errCh)
		for err := range errCh {
			require.NoError(t, err)
		}
	}
	for i := 0; i < 10; i++ {
		run()
	}
}

func TestClientApplicationServerGetDevices(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	u, err := url.Parse(device.Endpoints[0])
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	device.OwnershipStatus = grpcgwPb.Device_UNOWNED

	type args struct {
		req *pb.GetDevicesRequest
		srv pb.ClientApplication_GetDevicesServer
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
				srv: test.NewClientApplicationGetDevicesServer(ctx),
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
				srv: test.NewClientApplicationGetDevicesServer(ctx),
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
				srv: test.NewClientApplicationGetDevicesServer(ctx),
			},
			want: []*grpcgwPb.Device{
				device,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, teardown, err := test.NewClientApplicationServer(ctx)
			require.NoError(t, err)
			defer teardown()
			err = s.GetDevices(tt.args.req, tt.args.srv)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			server, ok := tt.args.srv.(*test.ClientApplicationGetDevicesServer)
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
