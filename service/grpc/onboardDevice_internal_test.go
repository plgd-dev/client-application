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

package grpc

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/schema"
	"github.com/plgd-dev/device/v2/schema/cloud"
	"github.com/stretchr/testify/require"
)

func TestOnboardInsecureDevice(t *testing.T) {
	// we don't have a insecure device simulator for this test, so we create a fake device
	// and try to onboard it
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	dev := device{Device: &core.Device{}}
	links := schema.ResourceLinks{
		{
			Href:          cloud.ResourceURI,
			ResourceTypes: []string{cloud.ResourceType},
		},
	}
	err := onboardInsecureDevice(ctx, &dev, links, &pb.OnboardDeviceRequest{
		DeviceId:                  "devId",
		CoapGatewayAddress:        "coaps+tcp://localhost:5684",
		AuthorizationCode:         "authCode",
		AuthorizationProviderName: "authProviderName",
		HubId:                     "hubId",
	})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "could not set cloud resource of device"))
}
