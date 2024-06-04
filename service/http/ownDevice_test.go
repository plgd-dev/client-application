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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/hub/v2/pkg/config/property/urischeme"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	hubTestOAuthServerService "github.com/plgd-dev/hub/v2/test/oauth-server/service"
	hubTestOAuthServerTest "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/plgd-dev/kit/v2/codec/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerOwnDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	getDevices(t, "")

	doSimpleOwn(t, dev.GetId(), http.StatusOK)

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResource, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.GetId()).ResourceHref("/light/1").Build()
	resp := httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	doDisown(t, dev.GetId(), http.StatusOK)
}

func createJwkKey(privateKey interface{}) (jwk.Key, error) {
	jwkKey, err := jwk.FromRaw(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwk: %w", err)
	}
	publicJwkKey, err := jwkKey.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	var publicKey any
	err = publicJwkKey.Raw(&publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	data, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal public key: %w", err)
	}
	if err = jwkKey.Set(jwk.KeyIDKey, uuid.NewSHA1(uuid.NameSpaceX500, data).String()); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jwk.KeyIDKey, err)
	}
	var alg string
	switch privateKey.(type) {
	case *rsa.PrivateKey:
		alg = jwa.RS256.String()
	case *ecdsa.PrivateKey:
		alg = jwa.ES256.String()
	default:
		alg = jwkKey.Algorithm().String()
	}
	if err = jwkKey.Set(jwk.AlgorithmKey, alg); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jwk.AlgorithmKey, err)
	}
	return jwkKey, nil
}

func MakeAccessToken(subject, audience, scopes string, validFor time.Duration) (jwt.Token, error) {
	token := jwt.New()

	if err := token.Set(jwt.SubjectKey, subject); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jwt.SubjectKey, err)
	}
	if err := token.Set(jwt.AudienceKey, audience); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jwt.AudienceKey, err)
	}
	now := time.Now()
	if err := token.Set(jwt.IssuedAtKey, now); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jwt.IssuedAtKey, err)
	}
	if validFor != 0 {
		if err := token.Set(jwt.ExpirationKey, now.Add(validFor)); err != nil {
			return nil, fmt.Errorf("failed to set %v: %w", jwt.ExpirationKey, err)
		}
	}
	if err := token.Set("scope", scopes); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", "scope", err)
	}
	return token, nil
}

func GetAccessTokenForUser(t *testing.T, user string) string {
	privKey, err := hubTestOAuthServerService.LoadPrivateKey(urischeme.URIScheme(os.Getenv("TEST_OAUTH_SERVER_ACCESS_TOKEN_PRIVATE_KEY")))
	require.NoError(t, err)
	jwkKey, err := createJwkKey(privKey)
	require.NoError(t, err)
	token, err := MakeAccessToken(user, "https://localhost:9080", "ocf.cloud", time.Hour)
	require.NoError(t, err)
	buf, err := json.Encode(token)
	require.NoError(t, err)
	jwtToken, err := MakeJWTPayload(privKey, jwkKey, buf)
	require.NoError(t, err)
	return string(jwtToken)
}

func MakeJWTPayload(key interface{}, jwkKey jwk.Key, data []byte) ([]byte, error) {
	hdr := jws.NewHeaders()
	if err := hdr.Set(jws.AlgorithmKey, jwkKey.Algorithm()); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jws.AlgorithmKey, err)
	}
	if err := hdr.Set(jws.TypeKey, `JWT`); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jws.TypeKey, err)
	}
	if err := hdr.Set(jws.KeyIDKey, jwkKey.KeyID()); err != nil {
		return nil, fmt.Errorf("failed to set %v: %w", jws.KeyIDKey, err)
	}
	payload, err := jws.Sign(data, jws.WithKey(jwkKey.Algorithm(), key), jws.WithHeaders(hdr))
	if err != nil {
		return nil, fmt.Errorf("failed to create UserToken: %w", err)
	}
	return payload, nil
}

func TestClientApplicationServerOwnDeviceRemoteProvisioning(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	tearDown := setupRemoteProvisioning(t)
	defer tearDown()

	initializeRemoteProvisioning(ctx, t)

	token := hubTestOAuthServerTest.GetDefaultAccessToken(t)
	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	getDevices(t, token)

	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, encodeToBody(t, &pb.OwnDeviceRequest{
		Timeout: (time.Second * 8).Nanoseconds(),
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp := httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCSRResp pb.OwnDeviceResponse
	decodeBody(t, resp.Body, &ownCSRResp)
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetState())

	certificate := signCSR(ctx, t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice+"/"+ownCSRResp.GetIdentityCertificateChallenge().GetState(), encodeToBody(t, &pb.FinishOwnDeviceRequest{
		Certificate: certificate,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCertificateResp pb.FinishOwnDeviceResponse
	decodeBody(t, resp.Body, &ownCertificateResp)

	// get resource with valid token
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResource, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).ResourceHref("/light/1").Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// get resource with token by another user
	user2token := GetAccessTokenForUser(t, "user2")
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResource, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(user2token).DeviceId(dev.GetId()).ResourceHref("/light/1").Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClientApplicationServerOwnDeviceRemoteProvisioningFails(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	tearDown := setupRemoteProvisioning(t)
	defer tearDown()

	initializeRemoteProvisioning(ctx, t)

	token := hubTestOAuthServerTest.GetDefaultAccessToken(t)
	ctx = kitNetGrpc.CtxWithToken(ctx, token)

	getDevices(t, token)

	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, encodeToBody(t, &pb.OwnDeviceRequest{
		Timeout: (time.Millisecond * 100).Nanoseconds(),
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp := httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCSRResp pb.OwnDeviceResponse
	decodeBody(t, resp.Body, &ownCSRResp)
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetState())

	// fail for timeout
	time.Sleep(time.Second)

	certificate := signCSR(ctx, t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice+"/"+ownCSRResp.GetIdentityCertificateChallenge().GetState(), encodeToBody(t, &pb.FinishOwnDeviceRequest{
		Certificate: certificate,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
