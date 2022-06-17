//***************************************************************************
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
//*************************************************************************

package http_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	hubTest "github.com/plgd-dev/hub/v2/test"
	"github.com/plgd-dev/kit/v2/codec/json"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerCreateDeleteResource(t *testing.T) {
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

	// create resource
	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DeviceResourceLink, bytes.NewBuffer(func() []byte {
		v, err := json.Encode(hubTest.MakeSwitchResourceDefaultData())
		require.NoError(t, err)
		return v
	}())).Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref(hubTest.TestResourceSwitchesHref).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// delete resource
	request = httpgwTest.NewRequest(http.MethodDelete, serviceHttp.DeviceResourceLink, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref(hubTest.TestResourceSwitchesInstanceHref("1")).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// duplicity delete
	request = httpgwTest.NewRequest(http.MethodDelete, serviceHttp.DeviceResourceLink, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref(hubTest.TestResourceSwitchesInstanceHref("1")).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
