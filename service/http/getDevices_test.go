package http_test

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	httpService "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetDevices(t *testing.T) {
	device := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	u, err := url.Parse(device.Endpoints[0])
	require.NoError(t, err)

	type args struct {
		accept       string
		useMulticast []string
		useEndpoints []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*pb.Device
	}{
		{
			name: "by multicast",
			args: args{
				useMulticast: []string{pb.GetDevicesRequest_IPV4.String()},
			},
			want: []*pb.Device{
				device,
			},
		},
		{
			name: "by ip",
			args: args{
				useEndpoints: []string{u.Hostname()},
			},
			want: []*pb.Device{
				device,
			},
		},
		{
			name: "by ip:port",
			args: args{
				useEndpoints: []string{u.Host},
			},
			want: []*pb.Device{
				device,
			},
		},
	}

	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.Connection.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httpgwTest.NewRequest(http.MethodGet, httpService.Devices, nil).
				Host(test.CLIENT_APPLICATIO_HTTP_HOST).Accept(tt.args.accept).AddQuery("useMulticast", tt.args.useMulticast...).AddQuery("useEndpoints", tt.args.useEndpoints...).AddQuery("timeout", "1000").Build()
			resp := httpgwTest.HTTPDo(t, request)
			defer func() {
				_ = resp.Body.Close()
			}()

			var got []*pb.Device
			for {
				var dev pb.Device
				err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
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
