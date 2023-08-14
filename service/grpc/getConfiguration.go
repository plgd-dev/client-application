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
	"time"

	"github.com/plgd-dev/client-application/pb"
)

func (s *ClientApplicationServer) GetConfiguration(ctx context.Context, _ *pb.GetConfigurationRequest) (*pb.GetConfigurationResponse, error) {
	info := s.info.Clone()
	devService := s.serviceDevice.Load()
	info.DeviceAuthenticationMode = pb.GetConfigurationResponse_UNINITIALIZED
	info.IsInitialized = false
	info.Owner = ""
	remoteProvisioning := s.GetConfig().RemoteProvisioning
	info.RemoteProvisioning = remoteProvisioning.Clone()
	if info.RemoteProvisioning == nil {
		info.RemoteProvisioning = &pb.RemoteProvisioning{}
	}
	info.RemoteProvisioning.CurrentTime = time.Now().UnixNano()
	if devService != nil {
		info.DeviceAuthenticationMode = devService.GetDeviceAuthenticationMode()
		info.IsInitialized = devService.IsInitialized()
		info.Owner = devService.GetOwner()
		if info.DeviceAuthenticationMode == pb.GetConfigurationResponse_X509 {
			info.RemoteProvisioning.Mode = pb.RemoteProvisioning_USER_AGENT
		} else {
			info.RemoteProvisioning.Mode = pb.RemoteProvisioning_MODE_NONE
		}
	}

	return info, nil
}
