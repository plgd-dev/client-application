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
	"io"
	"net/http"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/config/device"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
)

func doInitializePSK(t *testing.T, subjectID string, key string, expCode int) {
	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.Initialize, encodeToBody(t, &pb.InitializeRequest{
		PreSharedKey: &pb.InitializePreSharedKey{
			SubjectId: subjectID,
			Key:       key,
		},
	})).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp := httpgwTest.HTTPDo(t, request)
	require.Equal(t, expCode, resp.StatusCode)
	_ = resp.Body.Close()
}

func initializePSK(t *testing.T, subjectID string, key string) {
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

	doInitializePSK(t, subjectID, key, http.StatusOK)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.WellKnownConfiguration, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	configRespBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	var configResp1 pb.GetConfigurationResponse
	err = protojson.Unmarshal(configRespBody, &configResp1)
	require.NoError(t, err)
	require.True(t, configResp1.GetIsInitialized())
	require.Equal(t, configResp1.GetOwner(), subjectID)
	require.Equal(t, pb.GetConfigurationResponse_PRE_SHARED_KEY, configResp1.GetDeviceAuthenticationMode())
}

func setupClientApplicationForPSKInitialization(t *testing.T) func() {
	cfg := test.MakeConfig(t)
	cfg.Log.Level = zapcore.DebugLevel
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	cfg.Clients.Device.COAP.TLS.Authentication = device.AuthenticationUninitialized
	cfg.Clients.Device.COAP.TLS.PreSharedKey.Key = ""
	cfg.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr = ""
	shutDown := test.New(t, cfg)

	return shutDown
}

func doSimpleOwn(t *testing.T, deviceID string, expCode int) {
	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(deviceID).Build()
	resp := httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, expCode, resp.StatusCode)
}

func doDisown(t *testing.T, deviceID string, expCode int) {
	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(deviceID).Build()
	resp := httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, expCode, resp.StatusCode)
}

func TestClientApplicationServerInitializeWithPSK(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	shutDown := setupClientApplicationForPSKInitialization(t)
	defer shutDown()

	subjectID := "subjectID"
	initializePSK(t, subjectID, "key")
	doInitializePSK(t, subjectID, "key", http.StatusBadRequest)
	doSimpleOwn(t, dev.GetId(), http.StatusNotFound)
	getDevices(t, "")
	doSimpleOwn(t, dev.GetId(), http.StatusOK)

	// reset
	doReset(t, "", http.StatusOK)

	// initialize again
	initializePSK(t, subjectID, "key")
	doDisown(t, dev.GetId(), http.StatusNotFound)
	getDevices(t, "")
	doDisown(t, dev.GetId(), http.StatusOK)
}
