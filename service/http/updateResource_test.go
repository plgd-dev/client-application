package http_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	httpService "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerUpdateResource(t *testing.T) {
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
				accept:      httpService.ApplicationProtoJsonContentType,
				deviceID:    dev.Id,
				href:        doxm.ResourceURI,
				contentType: httpService.ApplicationJsonContentType,
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
				accept:      httpService.ApplicationProtoJsonContentType,
				deviceID:    dev.Id,
				href:        device.ResourceURI,
				contentType: httpService.ApplicationJsonContentType,
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
				contentType: httpService.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "unknown href",
			args: args{
				deviceID:    dev.Id,
				href:        "/unknown",
				contentType: httpService.ApplicationJsonContentType,
				body:        bytes.NewReader([]byte(`{"oxmsel":0}`)),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "forbidden - cannot establish TLS connection",
			args: args{
				deviceID:    dev.Id,
				href:        "/light/1",
				contentType: httpService.ApplicationJsonContentType,
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

	request := httpgwTest.NewRequest(http.MethodGet, httpService.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	for {
		var dev grpcgwPb.Device
		err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
	}
	_ = resp.Body.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodPut, httpService.DeviceResource, tt.args.body).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).DeviceId(tt.args.deviceID).ResourceHref(tt.args.href).ContentType(tt.args.contentType).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var got grpcgwPb.UpdateResourceResponse
			err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			hubTest.CheckProtobufs(t, tt.want, &got, hubTest.RequireToCheckFunc(require.Equal))
		})
	}
}
