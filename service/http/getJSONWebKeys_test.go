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
	"io"
	"net/http"
	"testing"

	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/test/config"
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetJSONWebKeys(t *testing.T) {
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.APIs.HTTP.UI.Enabled = true
	cfg.RemoteProvisioning.Mode = remoteProvisioning.Mode_UserAgent

	shutDown := test.New(t, cfg)
	defer shutDown()
	context.Background()

	oauthServerTearDown := hubTestOAuthServer.SetUp(t)
	defer oauthServerTearDown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(config.OAUTH_SERVER_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	jwks, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.WellKnownJWKs, bytes.NewReader(jwks)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	jwksClientApp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, jwks, jwksClientApp)
}
