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
	"context"
	"net/http"
	"testing"
	"time"

	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/kit/v2/codec/json"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetJSONWebKeys(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	tearDown := setupRemoteProvisioning(t)
	defer tearDown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	initializeRemoteProvisioning(ctx, t)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var jwksClientApp map[string]interface{}
	err := json.ReadFrom(resp.Body, &jwksClientApp)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, jwksClientApp)
}
