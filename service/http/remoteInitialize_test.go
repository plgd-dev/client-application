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
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/config/device"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	hubCAPb "github.com/plgd-dev/hub/v2/certificate-authority/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/plgd-dev/hub/v2/test/config"
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	hubTestService "github.com/plgd-dev/hub/v2/test/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
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

func initializeRemoteProvisioning(ctx context.Context, t *testing.T) {
	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var configResp pb.GetConfigurationResponse
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.False(t, configResp.GetIsInitialized())
	require.Equal(t, "", configResp.GetOwner())
	require.Equal(t, pb.GetConfigurationResponse_UNINITIALIZED, configResp.GetDeviceAuthenticationMode())
	require.Equal(t, pb.RemoteProvisioning_MODE_NONE, configResp.GetRemoteProvisioning().GetMode())

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(config.OAUTH_SERVER_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var jwks map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	require.NoError(t, err)

	pbJwks, err := structpb.NewStruct(jwks)
	require.NoError(t, err)

	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.Initialize, encodeToBody(t, &pb.InitializeRequest{
		Jwks: pbJwks,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var challengeResp pb.InitializeResponse
	decodeBody(t, resp.Body, &challengeResp)

	require.NotEmpty(t, challengeResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	require.NotEmpty(t, challengeResp.GetIdentityCertificateChallenge().GetState())
	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	certificate := signCSR(ctx, t, challengeResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.FinishInitialize(challengeResp.GetIdentityCertificateChallenge().GetState()), encodeToBody(t, &pb.FinishInitializeRequest{
		Certificate: certificate,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	configResp = pb.GetConfigurationResponse{}
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.True(t, configResp.GetIsInitialized())
	require.Equal(t, "1", configResp.GetOwner())
	require.Equal(t, pb.GetConfigurationResponse_X509, configResp.GetDeviceAuthenticationMode())
	require.Equal(t, pb.RemoteProvisioning_USER_AGENT, configResp.GetRemoteProvisioning().GetMode())

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCertificate, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func setupRemoteProvisioning(t *testing.T, services ...hubTestService.SetUpServicesConfig) func() {
	cfg := test.MakeConfig(t)
	cfg.Log.Level = zapcore.DebugLevel
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationUninitialized
	cfg.RemoteProvisioning.Mode = pb.RemoteProvisioning_MODE_NONE
	shutDown := test.New(t, cfg)

	var service hubTestService.SetUpServicesConfig
	if len(services) == 0 {
		service = hubTestService.SetUpServicesOAuth | hubTestService.SetUpServicesMachine2MachineOAuth | hubTestService.SetUpServicesCertificateAuthority
	} else {
		for _, s := range services {
			service |= s
		}
	}
	servicesShutdown := hubTestService.SetUpServices(context.TODO(), t, service)
	return func() {
		shutDown()
		servicesShutdown()
	}
}

func signCSR(ctx context.Context, t *testing.T, csr []byte) []byte {
	conn, err := grpc.NewClient(config.CERTIFICATE_AUTHORITY_HOST, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs:    hubTest.GetRootCertificatePool(t),
		MinVersion: tls.VersionTLS12,
	})))

	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()
	caClient := hubCAPb.NewCertificateAuthorityClient(conn)
	signResp, err := caClient.SignCertificate(ctx, &hubCAPb.SignCertificateRequest{
		CertificateSigningRequest: csr,
	})
	require.NoError(t, err)
	return signResp.GetCertificate()
}

func TestClientApplicationServerUpdateIdentityCertificate(t *testing.T) {
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationX509
	cfg.RemoteProvisioning.Mode = pb.RemoteProvisioning_USER_AGENT
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	shutDown := test.New(t, cfg)
	defer shutDown()

	servicesShutdown := hubTestService.SetUpServices(ctx, t, hubTestService.SetUpServicesOAuth|
		hubTestService.SetUpServicesMachine2MachineOAuth|hubTestService.SetUpServicesCertificateAuthority)
	defer servicesShutdown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var configResp pb.GetConfigurationResponse
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.False(t, configResp.GetIsInitialized())
	require.Equal(t, "", configResp.GetOwner())
	require.Equal(t, pb.GetConfigurationResponse_X509, configResp.GetDeviceAuthenticationMode())
	require.Equal(t, pb.RemoteProvisioning_USER_AGENT, configResp.GetRemoteProvisioning().GetMode())

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownJWKs, nil).Host(config.OAUTH_SERVER_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var jwks map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	require.NoError(t, err)

	pbJwks, err := structpb.NewStruct(jwks)
	require.NoError(t, err)

	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.Initialize, encodeToBody(t, &pb.InitializeRequest{
		Jwks: pbJwks,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var challengeResp pb.InitializeResponse
	decodeBody(t, resp.Body, &challengeResp)
	require.NotEmpty(t, challengeResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	require.NotEmpty(t, challengeResp.GetIdentityCertificateChallenge().GetState())

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCertificate, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	certificate := signCSR(ctx, t, challengeResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.FinishInitialize(challengeResp.GetIdentityCertificateChallenge().GetState()), encodeToBody(t, &pb.FinishInitializeRequest{
		Certificate: certificate,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.IdentityCertificate, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var clientCertificate pb.GetIdentityCertificateResponse
	decodeBody(t, resp.Body, &clientCertificate)
	require.Equal(t, certificate, clientCertificate.GetCertificate())

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = protojson.Unmarshal(configRespBody, &configResp)
	require.NoError(t, err)
	require.True(t, configResp.GetIsInitialized())
	require.Equal(t, "1", configResp.GetOwner())
}
