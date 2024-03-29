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

package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/client-application/test"
	"github.com/stretchr/testify/require"
)

func TestClientApplicationServerGetConfiguration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()

	d1, err := s.GetConfiguration(ctx, &pb.GetConfigurationRequest{})
	require.NoError(t, err)
	d1.RemoteProvisioning.CurrentTime = 0
	exp := test.NewServiceInformation()
	exp.IsInitialized = true
	exp.RemoteProvisioning = test.NewRemoteProvisioningConfig().Clone()
	require.Equal(t, exp, d1)
}

func TestClientApplicationServerGetConfigurationX509UserAgent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	cfg := test.MakeDeviceConfig()
	cfg.COAP.TLS.Authentication = device.AuthenticationX509
	remoteProvisioningConfig := test.NewRemoteProvisioningConfig()
	remoteProvisioningConfig.Mode = pb.RemoteProvisioning_USER_AGENT
	s, teardown, err := test.NewClientApplicationServer(ctx, test.WithDeviceConfig(cfg), test.WithRemoteProvisioningConfig(remoteProvisioningConfig))
	require.NoError(t, err)
	defer teardown()

	d1, err := s.GetConfiguration(ctx, &pb.GetConfigurationRequest{})
	require.NoError(t, err)
	d1.RemoteProvisioning.CurrentTime = 0
	exp := test.NewServiceInformation()
	exp.IsInitialized = false
	exp.Owner = ""
	exp.DeviceAuthenticationMode = pb.GetConfigurationResponse_X509
	exp.RemoteProvisioning = remoteProvisioningConfig.Clone()
	require.Equal(t, exp, d1)
}
