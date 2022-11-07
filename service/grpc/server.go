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

	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/plgd-dev/client-application/pb"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/go-coap/v3/pkg/sync"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientApplicationServer struct {
	pb.UnimplementedClientApplicationServer

	serviceDevice            *serviceDevice.Service
	info                     *ServiceInformation
	logger                   log.Logger
	devices                  *sync.Map[uuid.UUID, *device]
	csrCache                 *ttlcache.Cache[uuid.UUID, bool]
	remoteProvisioningConfig remoteProvisioning.Config
	jwksCache                atomic.Pointer[JSONWebKeyCache]
	remoteOwnSignCache       *sync.Map[uuid.UUID, *remoteSign]
}

func NewClientApplicationServer(remoteProvisioningConfig remoteProvisioning.Config, serviceDevice *serviceDevice.Service, info *ServiceInformation, logger log.Logger) *ClientApplicationServer {
	csrCache := ttlcache.New[uuid.UUID, bool]()
	go csrCache.Start()
	return &ClientApplicationServer{
		serviceDevice:            serviceDevice,
		info:                     info,
		logger:                   logger,
		csrCache:                 csrCache,
		remoteProvisioningConfig: remoteProvisioningConfig,
		remoteOwnSignCache:       sync.NewMap[uuid.UUID, *remoteSign](),
		devices:                  sync.NewMap[uuid.UUID, *device](),
	}
}

func (s *ClientApplicationServer) Version() string {
	return s.info.Version
}

func (s *ClientApplicationServer) Close() {
	s.csrCache.Stop()
}

func (s *ClientApplicationServer) getDevice(deviceID uuid.UUID) (*device, error) {
	dev, ok := s.devices.Load(deviceID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "device %v not found", deviceID)
	}
	return dev, nil
}

func (s *ClientApplicationServer) deleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	dev, ok := s.devices.LoadAndDelete(deviceID)
	if !ok {
		return nil
	}
	if err := dev.Close(ctx); err != nil {
		return status.Errorf(codes.Internal, "cannot close device %v connections: %v", deviceID, err)
	}
	return nil
}
