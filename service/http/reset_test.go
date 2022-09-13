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
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerReset(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cancel()

	tearDown := setupRemoteProvisioning(t)
	defer tearDown()

	initializeRemoteProvisioning(ctx, t)

	token := hubTestOAuthServer.GetDefaultAccessToken(t)

	// reset
	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.Reset, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// try again reset - should fail
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.Reset, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	initializeRemoteProvisioning(ctx, t)
}