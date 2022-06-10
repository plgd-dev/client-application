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
	"github.com/plgd-dev/device/schema/configuration"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	httpgwTest "github.com/plgd-dev/hub/v2/http-gateway/test"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGatewayServerGetDevice(t *testing.T) {
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

	newName := test.DevsimName + "_new"

	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.DeviceResource, bytes.NewBufferString(`{"n":"`+newName+`"}`)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref(configuration.ResourceURI).ContentType(serviceHttp.ApplicationJsonContentType).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.Device, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(serviceHttp.ApplicationProtoJsonContentType).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var device grpcgwPb.Device
	err := httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &device)
	require.NoError(t, err)

	var v plgdDevice.Device
	err = cbor.Decode(device.GetData().GetContent().GetData(), &v)
	require.NoError(t, err)
	require.Equal(t, newName, v.Name)

	request = httpgwTest.NewRequest(http.MethodGet, serviceHttp.Device, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).Accept(serviceHttp.ApplicationProtoJsonContentType).DeviceId(dev.Id).AddQuery(serviceHttp.UseCacheQueryKey).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var device1 grpcgwPb.Device
	err = httpgwTest.Unmarshal(resp.StatusCode, resp.Body, &device1)
	require.NoError(t, err)

	var v1 plgdDevice.Device
	err = cbor.Decode(device1.GetData().GetContent().GetData(), &v1)
	require.NoError(t, err)
	require.Equal(t, newName, v1.Name)

	request = httpgwTest.NewRequest(http.MethodPut, serviceHttp.DeviceResource, bytes.NewBufferString(`{"n":"`+test.DevsimName+`"}`)).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).ResourceHref(configuration.ResourceURI).ContentType(serviceHttp.ApplicationJsonContentType).Build()
	resp = httpgwTest.HTTPDo(t, request)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	request = httpgwTest.NewRequest(http.MethodPost, serviceHttp.DisownDevice, nil).
		Host(test.CLIENT_APPLICATION_HTTP_HOST).DeviceId(dev.Id).Build()
	resp = httpgwTest.HTTPDo(t, request)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
