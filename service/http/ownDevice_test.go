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
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	hubTestOAuthServer "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerOwnDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	for {
		var dev grpcgwPb.Device
		err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &dev)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
	}
	_ = resp.Body.Close()

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.DeviceResource, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref("/light/1").Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClientApplicationServerOwnDeviceRemoteProvisioning(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cancel()
	tearDown := setupRemoteProvisioning(ctx, t)
	defer tearDown()

	token := hubTestOAuthServer.GetDefaultAccessToken(t)
	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	request := httpgwTest.NewRequest(http.MethodGet, serviceHttp.Devices, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, encodeToBody(t, &pb.OwnDeviceRequest{
		GetIdentityCsr: &pb.OwnDeviceRequest_GetIdentityCsr{
			Timeout: (time.Second * 3600).Nanoseconds(),
		},
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCSRResp pb.OwnDeviceResponse
	decodeBody(t, resp.Body, &ownCSRResp)
	require.NotEmpty(t, ownCSRResp.GetGetIdentityCsr())
	require.NotEmpty(t, ownCSRResp.GetGetIdentityCsr().GetState())

	certificate := signCSR(ctx, t, ownCSRResp.GetGetIdentityCsr().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, encodeToBody(t, &pb.OwnDeviceRequest{
		SetIdentityCertificate: &pb.UpdateIdentityCertificateRequest{
			Certificate: certificate,
			State:       ownCSRResp.GetGetIdentityCsr().GetState(),
		},
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.Id).Build()
	time.Sleep(time.Second)
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCertificateResp pb.OwnDeviceResponse
	decodeBody(t, resp.Body, &ownCertificateResp)
	require.Empty(t, ownCertificateResp.GetGetIdentityCsr())
	require.Empty(t, ownCertificateResp.GetGetIdentityCsr().GetState())

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
