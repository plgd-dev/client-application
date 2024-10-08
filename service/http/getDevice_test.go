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
	"net/http"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/v2/schema/configuration"
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	pkgHttpPb "github.com/plgd-dev/hub/v2/pkg/net/http/pb"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetDevice(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	cfg := test.MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	shutDown := test.New(t, cfg)
	defer shutDown()

	getDevices(t, "")

	request := httpgwTest.NewRequest(http.MethodPost, serviceHttp.OwnDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.GetId()).Build()
	resp := httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	newName := test.DevsimName + "_new"

	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.DeviceResource, bytes.NewBufferString(`{"n":"`+newName+`"}`)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.GetId()).ResourceHref(configuration.ResourceURI).ContentType(serviceHttp.ApplicationJsonContentType).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.Device, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(serviceHttp.ApplicationProtoJsonContentType).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var device grpcgwPb.Device
	err := pkgHttpPb.Unmarshal(resp.StatusCode, resp.Body, &device)
	require.NoError(t, err)

	var v plgdDevice.Device
	err = cbor.Decode(device.GetData().GetContent().GetData(), &v)
	require.NoError(t, err)
	require.Equal(t, newName, v.Name)

	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.DeviceResource, bytes.NewBufferString(`{"n":"`+test.DevsimName+`"}`)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.GetId()).ResourceHref(configuration.ResourceURI).ContentType(serviceHttp.ApplicationJsonContentType).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.GetId()).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
