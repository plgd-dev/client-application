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

package http_test

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	httpTest "github.com/plgd-dev/hub/v2/test/http"
	"github.com/stretchr/testify/require"
)

func getDevices(t *testing.T, token string) {
	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST)
	if token != "" {
		request = request.AuthToken(token)
	}
	resp := httpgwTest.HTTPDo(t, request.Build())
	require.Equal(t, http.StatusOK, resp.StatusCode)
	for {
		var dev grpcgwPb.Device
		err := httpTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
	}
	_ = resp.Body.Close()
}

func TestClientApplicationServerGetDevices(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	u, err := url.Parse(device.GetEndpoints()[0])
	require.NoError(t, err)

	type args struct {
		accept       string
		useMulticast []string
		useEndpoints []string
		timeout      time.Duration
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
				useMulticast: []string{pb.GetDevicesRequest_IPV4.String()},
			},
			want: []*grpcgwPb.Device{
				device,
			},
		},
		{
			name: "by ip",
			args: args{
				useEndpoints: []string{u.Hostname()},
			},
			want: []*grpcgwPb.Device{
				device,
			},
		},
		{
			name: "by ip:port",
			args: args{
				useEndpoints: []string{u.Host},
			},
			want: []*grpcgwPb.Device{
				device,
			},
		},
	}

	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).AddQuery(serviceHttp.UseMulticastQueryKey, tt.args.useMulticast...).AddQuery(serviceHttp.UseEndpointsQueryKey, tt.args.useEndpoints...).AddQuery(serviceHttp.TimeoutQueryKey, strconv.FormatInt(int64(tt.args.timeout/time.Millisecond), 10)).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()

			var got []*grpcgwPb.Device
			for {
				var dev grpcgwPb.Device
				err := httpTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				got = append(got, &dev)
			}
			require.Equal(t, len(tt.want), len(got))
		})
	}
}
