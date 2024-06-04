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
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/v2/schema/configuration"
	"github.com/plgd-dev/device/v2/schema/device"
	"github.com/plgd-dev/device/v2/schema/doxm"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	hubTest "github.com/plgd-dev/hub/v2/test"
	httpTest "github.com/plgd-dev/hub/v2/test/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerUpdateResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	dev.Data.OpenTelemetryCarrier = map[string]string{}

	type args struct {
		accept      string
		deviceID    string
		href        string
		contentType string
		body        io.Reader
	}
	tests := []struct {
		name     string
		args     args
		want     *grpcgwPb.UpdateResourceResponse
		wantErr  bool
		wantCode int
	}{
		{
			name: "doxm update",
			args: args{
				accept:      serviceHttp.ApplicationProtoJsonContentType,
				deviceID:    dev.GetId(),
				href:        doxm.ResourceURI,
				contentType: serviceHttp.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			want: &grpcgwPb.UpdateResourceResponse{
				Data: &events.ResourceUpdated{
					Content:              &commands.Content{},
					Status:               commands.Status_OK,
					OpenTelemetryCarrier: map[string]string{},
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "device resource",
			args: args{
				accept:      serviceHttp.ApplicationProtoJsonContentType,
				deviceID:    dev.GetId(),
				href:        device.ResourceURI,
				contentType: serviceHttp.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "unknown device",
			args: args{
				deviceID:    uuid.NewString(),
				href:        device.ResourceURI,
				contentType: serviceHttp.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "unknown href",
			args: args{
				deviceID:    dev.GetId(),
				href:        "/unknown",
				contentType: serviceHttp.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "forbidden - cannot establish TLS connection",
			args: args{
				deviceID:    dev.GetId(),
				href:        configuration.ResourceURI,
				contentType: serviceHttp.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
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
			request := httpgwTest.NewRequest(http.MethodPut, serviceHttp.DeviceResource, tt.args.body).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).DeviceId(tt.args.deviceID).ResourceHref(tt.args.href).ContentType(tt.args.contentType).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var got grpcgwPb.UpdateResourceResponse
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
