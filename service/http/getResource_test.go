package http_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	httpService "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/device"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})

	type args struct {
		accept   string
		deviceID string
		href     string
	}
	tests := []struct {
		name     string
		args     args
		want     *pb.GetResourceResponse
		wantErr  bool
		wantCode int
	}{
		{
			name: "device resource",
			args: args{
				accept:   httpService.ApplicationProtoJsonContentType,
				deviceID: dev.Id,
				href:     device.ResourceURI,
			},
			want: &pb.GetResourceResponse{
				Content: dev.Content,
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
				deviceID: dev.Id,
				href:     "/unknown",
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "unavailable - cannot establish TLS connection",
			args: args{
				deviceID: dev.Id,
				href:     "/light/1",
			},
			wantErr:  true,
			wantCode: http.StatusServiceUnavailable,
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
		var dev pb.Device
		err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
	}
	_ = resp.Body.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodGet, httpService.DeviceResource, nil).
				Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(tt.args.accept).DeviceId(tt.args.deviceID).ResourceHref(tt.args.href).Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var got pb.GetResourceResponse
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
