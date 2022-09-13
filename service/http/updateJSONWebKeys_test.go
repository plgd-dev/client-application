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
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/device"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/test/config"
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/plgd-dev/kit/v2/codec/json"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerUpdateJSONWebKeys(t *testing.T) {
	// need to wait for the device
	_ = test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})

	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.RemoteProvisioning.Mode = remoteProvisioning.Mode_UserAgent
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationX509

	shutDown := test.New(t, cfg)
	defer shutDown()
	context.Background()

	oauthServerTearDown := hubTestOAuthServer.SetUp(t)
	defer oauthServerTearDown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(config.OAUTH_SERVER_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	jwksBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// update without token
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.WellKnownJWKs, bytes.NewReader(jwksBody)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// update with invalid token
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.WellKnownJWKs, bytes.NewReader(jwksBody)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken("invalidToken").Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// update with token
	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.WellKnownJWKs, bytes.NewReader(jwksBody)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// refresh jwks with token
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.WellKnownJWKs, bytes.NewReader(jwksBody)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// get jwks
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	var jwksClientApp map[string]interface{}
	err = json.ReadFrom(resp.Body, &jwksClientApp)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var jwks map[string]interface{}
	err = json.Decode(jwksBody, &jwks)
	require.NoError(t, err)
	require.Equal(t, jwks, jwksClientApp)

	// get devices without token
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// get devices with token
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	foundDevices := false
	for {
		var dev grpcgwPb.Device
		err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		foundDevices = true
	}
	_ = resp.Body.Close()
	require.True(t, foundDevices)
}
