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
	"crypto/tls"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/device"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/client-application/test"
	hubCAPb "github.com/plgd-dev/hub/v2/certificate-authority/pb"
	caService "github.com/plgd-dev/hub/v2/certificate-authority/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/plgd-dev/hub/v2/test/config"
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestClientApplicationServerUpdateIdentityCertificate(t *testing.T) {
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationX509
	// /cfg.APIs.HTTP.UI.Enabled = true
	cfg.RemoteProvisioning.Mode = remoteProvisioning.Mode_UserAgent
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	shutDown := test.New(t, cfg)
	defer shutDown()
	context.Background()

	oauthServerTearDown := hubTestOAuthServer.SetUp(t)
	defer oauthServerTearDown()

	caShutdown := caService.SetUp(t)
	defer caShutdown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(config.OAUTH_SERVER_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	jwks, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.WellKnownJWKs, bytes.NewReader(jwks)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCsr, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	csrBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var csr pb.GetIdentityCSRResponse
	err = protojson.Unmarshal(csrBody, &csr)
	require.NoError(t, err)

	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	conn, err := grpc.Dial(config.CERTIFICATE_AUTHORITY_HOST, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs: hubTest.GetRootCertificatePool(t),
	})))
	require.NoError(t, err)
	caClient := hubCAPb.NewCertificateAuthorityClient(conn)
	signResp, err := caClient.SignCertificate(ctx, &hubCAPb.SignCertificateRequest{
		CertificateSigningRequest: []byte(csr.CertificateSigningRequest),
	})
	require.NoError(t, err)
	updateCertificate := pb.UpdateIdentityCertificateRequest{
		Certificate: string(signResp.Certificate),
		State:       csr.State,
	}
	updateCertificateBody, err := protojson.Marshal(&updateCertificate)
	require.NoError(t, err)

	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.IdentityCertificate, io.NopCloser(bytes.NewReader(updateCertificateBody))).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCertificate, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	clientCertificateBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var clientCertificate pb.GetIdentityCertificateResponse
	err = protojson.Unmarshal(clientCertificateBody, &clientCertificate)
	require.NoError(t, err)
	require.Equal(t, string(signResp.Certificate), clientCertificate.Certificate)
}
