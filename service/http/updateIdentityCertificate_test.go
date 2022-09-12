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
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func decodeBody(t *testing.T, r io.Reader, v protoreflect.ProtoMessage) {
	csrBody, err := io.ReadAll(r)
	require.NoError(t, err)
	err = protojson.Unmarshal(csrBody, v)
	require.NoError(t, err)
}

func encodeToBody(t *testing.T, v protoreflect.ProtoMessage) io.Reader {
	body, err := protojson.Marshal(v)
	require.NoError(t, err)
	return io.NopCloser(bytes.NewReader(body))
}

func setupRemoteProvisioning(ctx context.Context, t *testing.T) func() {
	cfg := test.MakeConfig(t)
	cfg.Log.Level = zapcore.DebugLevel
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationX509
	cfg.RemoteProvisioning.Mode = remoteProvisioning.Mode_UserAgent
	shutDown := test.New(t, cfg)
	oauthServerTearDown := hubTestOAuthServer.SetUp(t)
	caShutdown := caService.SetUp(t)

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var configResp pb.GetConfigurationResponse
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.False(t, configResp.GetIsInitialized())
	require.Equal(t, configResp.GetDeviceAuthenticationMode(), pb.GetConfigurationResponse_X509)
	require.Equal(t, configResp.GetRemoteProvisioning().GetMode(), pb.RemoteProvisioning_USER_AGENT)

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

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCsr, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var csr pb.GetIdentityCSRResponse
	decodeBody(t, resp.Body, &csr)

	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	conn, err := grpc.Dial(config.CERTIFICATE_AUTHORITY_HOST, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs:    hubTest.GetRootCertificatePool(t),
		MinVersion: tls.VersionTLS12,
	})))
	require.NoError(t, err)
	caClient := hubCAPb.NewCertificateAuthorityClient(conn)
	signResp, err := caClient.SignCertificate(ctx, &hubCAPb.SignCertificateRequest{
		CertificateSigningRequest: []byte(csr.CertificateSigningRequest),
	})
	require.NoError(t, err)
	initializeReq := pb.InitializeRequest{
		X509: &pb.UpdateIdentityCertificateRequest{
			Certificate: string(signResp.Certificate),
			State:       csr.State,
		},
	}

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.Initialize, encodeToBody(t, &initializeReq)).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return func() {
		shutDown()
		oauthServerTearDown()
		caShutdown()
	}
}

func signCSR(ctx context.Context, t *testing.T, csr string) string {
	conn, err := grpc.Dial(config.CERTIFICATE_AUTHORITY_HOST, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs:    hubTest.GetRootCertificatePool(t),
		MinVersion: tls.VersionTLS12,
	})))

	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()
	caClient := hubCAPb.NewCertificateAuthorityClient(conn)
	signResp, err := caClient.SignCertificate(ctx, &hubCAPb.SignCertificateRequest{
		CertificateSigningRequest: []byte(csr),
	})
	require.NoError(t, err)
	return string(signResp.Certificate)
}

func TestClientApplicationServerUpdateIdentityCertificate(t *testing.T) {
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationX509
	cfg.RemoteProvisioning.Mode = remoteProvisioning.Mode_UserAgent
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	shutDown := test.New(t, cfg)
	defer shutDown()

	oauthServerTearDown := hubTestOAuthServer.SetUp(t)
	defer oauthServerTearDown()

	caShutdown := caService.SetUp(t)
	defer caShutdown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var configResp pb.GetConfigurationResponse
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.False(t, configResp.GetIsInitialized())
	require.Equal(t, configResp.GetDeviceAuthenticationMode(), pb.GetConfigurationResponse_X509)
	require.Equal(t, configResp.GetRemoteProvisioning().GetMode(), pb.RemoteProvisioning_USER_AGENT)

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

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCsr, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var csr pb.GetIdentityCSRResponse
	decodeBody(t, resp.Body, &csr)

	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	certificate := signCSR(ctx, t, csr.CertificateSigningRequest)

	initializeReq := pb.InitializeRequest{
		X509: &pb.UpdateIdentityCertificateRequest{
			Certificate: certificate,
			State:       csr.State,
		},
	}

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.Initialize, encodeToBody(t, &initializeReq)).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.IdentityCertificate, encodeToBody(t, initializeReq.GetX509())).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCertificate, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var clientCertificate pb.GetIdentityCertificateResponse
	decodeBody(t, resp.Body, &clientCertificate)
	require.Equal(t, certificate, clientCertificate.Certificate)
}
