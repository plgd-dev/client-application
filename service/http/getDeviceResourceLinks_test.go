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
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	hubTest "github.com/plgd-dev/hub/v2/test"
	hubTestPb "github.com/plgd-dev/hub/v2/test/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetDeviceResourceLinks(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})

	type args struct {
		deviceID string
		accept   string
	}
	tests := []struct {
		name     string
		args     args
		want     *events.ResourceLinksPublished
		wantErr  bool
		wantCode int
	}{
		{
			name: "device links",
			args: args{
				deviceID: dev.GetId(),
			},
			want: hubTestPb.CleanUpResourceLinksPublished(&events.ResourceLinksPublished{
				DeviceId:  dev.GetId(),
				Resources: commands.SchemaResourceLinksToResources(test.GetDeviceResourceLinks(), time.Time{}),
			}, true),
			wantCode: http.StatusOK,
		},
		{
			name: "unknown device",
			args: args{
				deviceID: uuid.NewString(),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
	}

	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	getDevices(t, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResourceLinks, nil).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).DeviceId(tt.args.deviceID).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var val events.ResourceLinksPublished
			err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &val)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			got := hubTestPb.CleanUpResourceLinksPublished(&val, true)
			for _, r := range got.GetResources() {
				require.NotEmpty(t, r.GetEndpointInformations())
				r.EndpointInformations = nil
				require.NotEmpty(t, r.GetPolicy())
				r.Policy = nil
				require.NotEmpty(t, r.GetAnchor())
				r.Anchor = ""
			}
			require.NoError(t, err)
			hubTest.CheckProtobufs(t, tt.want, &got, hubTest.RequireToCheckFunc(require.Equal))
		})
	}
}
