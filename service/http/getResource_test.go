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
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/v2/schema/configuration"
	"github.com/plgd-dev/device/v2/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	hubTest "github.com/plgd-dev/hub/v2/test"
	httpTest "github.com/plgd-dev/hub/v2/test/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	dev.Data.OpenTelemetryCarrier = map[string]string{}

	type args struct {
		accept   string
		deviceID string
		href     string
	}
	tests := []struct {
		name     string
		args     args
		want     *grpcgwPb.Resource
		wantErr  bool
		wantCode int
	}{
		{
			name: "device resource",
			args: args{
				accept:   serviceHttp.ApplicationProtoJsonContentType,
				deviceID: dev.GetId(),
				href:     device.ResourceURI,
			},
			want: &grpcgwPb.Resource{
				Data:  dev.GetData(),
				Types: []string{"oic.d.cloudDevice", "oic.wk.d"},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "unknown device",
			args: args{
				deviceID: uuid.NewString(),
				href:     device.ResourceURI,
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "unknown href",
			args: args{
				deviceID: dev.GetId(),
				href:     "/unknown",
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "forbidden - cannot establish TLS connection",
			args: args{
				deviceID: dev.GetId(),
				href:     configuration.ResourceURI,
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
	}

	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	getDevices(t, "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResource, nil).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).DeviceId(tt.args.deviceID).ResourceHref(tt.args.href).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var got grpcgwPb.Resource
			err := httpTest.Unmarshal(resp.StatusCode, resp.Body, &got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			hubTest.CheckProtobufs(t, tt.want, &got, hubTest.RequireToCheckFunc(require.Equal))
		})
	}
}
