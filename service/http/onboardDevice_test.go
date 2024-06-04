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
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	hubGrpcGwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/plgd-dev/hub/v2/test/config"
	"github.com/plgd-dev/hub/v2/test/device/ocf"
	httpTest "github.com/plgd-dev/hub/v2/test/http"
	hubTestOAuthServerTest "github.com/plgd-dev/hub/v2/test/oauth-server/test"
	hubTestService "github.com/plgd-dev/hub/v2/test/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestClientApplicationServerOnboardDeviceRemoteProvisioning(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	tearDown := setupRemoteProvisioning(t, hubTestService.SetUpServicesOAuth|
		hubTestService.SetUpServicesCertificateAuthority|
		hubTestService.SetUpServicesId|
		hubTestService.SetUpServicesCertificateAuthority|
		hubTestService.SetUpServicesCoapGateway|
		hubTestService.SetUpServicesGrpcGateway|
		hubTestService.SetUpServicesResourceAggregate|
		hubTestService.SetUpServicesResourceDirectory,
	)
	defer tearDown()

	initializeRemoteProvisioning(ctx, t)

	token := hubTestOAuthServerTest.GetDefaultAccessToken(t)
	ctx = kitNetGrpc.CtxWithToken(ctx, token)
	getDevices(t, token)

	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, encodeToBody(t, &pb.OwnDeviceRequest{
		Timeout: (time.Second * 8).Nanoseconds(),
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCSRResp pb.OwnDeviceResponse
	decodeBody(t, resp.Body, &ownCSRResp)
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	require.NotEmpty(t, ownCSRResp.GetIdentityCertificateChallenge().GetState())

	// own
	certificate := signCSR(ctx, t, ownCSRResp.GetIdentityCertificateChallenge().GetCertificateSigningRequest())
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice+"/"+ownCSRResp.GetIdentityCertificateChallenge().GetState(), encodeToBody(t, &pb.FinishOwnDeviceRequest{
		Certificate: certificate,
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var ownCertificateResp pb.FinishOwnDeviceResponse
	decodeBody(t, resp.Body, &ownCertificateResp)

	cloudConn, err := grpc.NewClient(config.GRPC_GW_HOST, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs:    hubTest.GetRootCertificatePool(t),
		MinVersion: tls.VersionTLS12,
	})))
	require.NoError(t, err)
	defer func() {
		err = cloudConn.Close()
		require.NoError(t, err)
	}()
	c := hubGrpcGwPb.NewGrpcGatewayClient(cloudConn)
	subClient, err := c.SubscribeToEvents(ctx)
	require.NoError(t, err)
	err = subClient.Send(&hubGrpcGwPb.SubscribeToEvents{
		CorrelationId: "allEvents",
		Action: &hubGrpcGwPb.SubscribeToEvents_CreateSubscription_{
			CreateSubscription: &hubGrpcGwPb.SubscribeToEvents_CreateSubscription{},
		},
	})
	require.NoError(t, err)
	defer func() {
		err = subClient.CloseSend()
		require.NoError(t, err)
	}()
	ev, err := subClient.Recv()
	require.NoError(t, err)
	expectedEvent := &hubGrpcGwPb.Event{
		SubscriptionId: ev.GetSubscriptionId(),
		CorrelationId:  "allEvents",
		Type: &hubGrpcGwPb.Event_OperationProcessed_{
			OperationProcessed: &hubGrpcGwPb.Event_OperationProcessed{
				ErrorStatus: &hubGrpcGwPb.Event_OperationProcessed_ErrorStatus{
					Code: hubGrpcGwPb.Event_OperationProcessed_ErrorStatus_OK,
				},
			},
		},
	}
	hubTest.CheckProtobufs(t, expectedEvent, ev, hubTest.RequireToCheckFunc(require.Equal))

	// get configuration
	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var httpClientCfg pb.GetConfigurationResponse
	err = httpTest.Unmarshal(resp.StatusCode, resp.Body, &httpClientCfg)
	require.NoError(t, err)

	// oauth server url to host
	oauthServerUrl := httpClientCfg.GetRemoteProvisioning().GetAuthority()
	oauthServer, err := url.Parse(oauthServerUrl)
	require.NoError(t, err)
	require.Equal(t, config.OAUTH_SERVER_HOST, oauthServer.Host)

	// onboard
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OnboardDevice, encodeToBody(t, &pb.OnboardDeviceRequest{
		CoapGatewayAddress:        httpClientCfg.GetRemoteProvisioning().GetCoapGateway(),
		AuthorizationCode:         hubTestOAuthServerTest.GetAuthorizationCode(t, oauthServer.Host, httpClientCfg.GetRemoteProvisioning().GetDeviceOauthClient().GetClientId(), dev.GetId(), ""),
		HubId:                     httpClientCfg.GetRemoteProvisioning().GetId(),
		AuthorizationProviderName: httpClientCfg.GetRemoteProvisioning().GetDeviceOauthClient().GetProviderName(),
		CertificateAuthorities:    httpClientCfg.GetRemoteProvisioning().GetCertificateAuthorities(),
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	hubTest.WaitForDevice(t, subClient, ocf.NewDevice(dev.GetId(), test.DevsimName), expectedEvent.GetSubscriptionId(), expectedEvent.GetCorrelationId(), hubTest.GetAllBackendResourceLinks())

	// offboard
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.OffboardDevice, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// disown
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).AuthToken(token).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
