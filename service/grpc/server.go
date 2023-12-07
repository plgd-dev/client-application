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
	"sync"

	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/config"
	configGrpc "github.com/plgd-dev/client-application/service/config/grpc"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	coapSync "github.com/plgd-dev/go-coap/v3/pkg/sync"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientApplicationServer struct {
	pb.UnimplementedClientApplicationServer

	serviceDevice      atomic.Pointer[serviceDevice.Service]
	info               *pb.GetConfigurationResponse
	logger             log.Logger
	devices            *coapSync.Map[uuid.UUID, *device]
	csrCache           *ttlcache.Cache[uuid.UUID, *serviceDevice.Service]
	config             *atomic.Pointer[config.Config]
	jwksCache          atomic.Pointer[JSONWebKeyCache]
	remoteOwnSignCache *coapSync.Map[uuid.UUID, *remoteSign]

	initializationMutex sync.Mutex
}

func NewClientApplicationServer(cfg *atomic.Pointer[config.Config], devService *serviceDevice.Service, info *configGrpc.ServiceInformation, logger log.Logger) *ClientApplicationServer {
	csrCache := ttlcache.New[uuid.UUID, *serviceDevice.Service]()
	go csrCache.Start()
	s := ClientApplicationServer{
		info:               pb.NewGetConfigurationResponse(info),
		logger:             logger,
		csrCache:           csrCache,
		config:             cfg,
		remoteOwnSignCache: coapSync.NewMap[uuid.UUID, *remoteSign](),
		devices:            coapSync.NewMap[uuid.UUID, *device](),
	}
	if devService != nil {
		s.init(context.Background(), devService)
	}
	return &s
}

func (s *ClientApplicationServer) Version() string {
	return s.info.Version
}

func (s *ClientApplicationServer) GetConfig() config.Config {
	cfg := s.config.Load()
	return *cfg
}

func (s *ClientApplicationServer) StoreConfig(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid configuration: %v", err)
	}
	if err := cfg.Store(); err != nil {
		return status.Errorf(codes.Internal, "cannot store configuration: %v", err)
	}
	s.config.Store(cfg)
	return nil
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
