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
	"github.com/plgd-dev/client-application/pb"
	configDevice "github.com/plgd-dev/client-application/service/config/device"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errAlreadyInitialized = "already initialized"

func (s *ClientApplicationServer) InitializeRemoteProvisioning(ctx context.Context, req *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	err := s.UpdateJSONWebKeys(ctx, req.GetJwks())
	if err != nil {
		return nil, err
	}
	cfg := s.GetConfig()
	cfg.Clients.Device.COAP.TLS.Authentication = configDevice.AuthenticationX509
	devService, err := serviceDevice.New(context.Background(), func() configDevice.Config {
		return cfg.Clients.Device
	}, s.logger)
	if err != nil {
		return nil, err
	}

	respCsr, err := s.getIdentityCSR(ctx, devService)
	if err != nil {
		if err2 := s.reset(ctx); err2 != nil {
			s.logger.Warnf("cannot reset when initialize remote provisioning fails(%v): %w", err, err2)
		}
		return nil, err
	}
	return &pb.InitializeResponse{
		IdentityCertificateChallenge: respCsr,
	}, nil
}

func (s *ClientApplicationServer) init(ctx context.Context, devService *serviceDevice.Service) {
	err := s.reset(ctx)
	if err != nil {
		s.logger.Warnf("cannot reset during initialize fails: %v", err)
	}
	s.serviceDevice.Store(devService)
	go func() {
		err := devService.Serve()
		if err != nil {
			s.logger.Warnf("device service cannot serve coap connections: %v", err)
		}
	}()
}

func (s *ClientApplicationServer) UpdatePSK(subjectUUID, key string) error {
	cfg := s.GetConfig()
	if cfg.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr == subjectUUID && cfg.Clients.Device.COAP.TLS.PreSharedKey.Key == key {
		return nil
	}
	cfg.Clients.Device.COAP.TLS.PreSharedKey.Key = key
	cfg.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr = subjectUUID
	return s.StoreConfig(cfg)
}

func (s *ClientApplicationServer) initWithPSK(ctx context.Context, subjectUUID, key string) error {
	err := s.UpdatePSK(subjectUUID, key)
	if err != nil {
		return err
	}
	cfg := s.GetConfig()
	cfg.Clients.Device.COAP.TLS.Authentication = configDevice.AuthenticationPreSharedKey
	cfg.RemoteProvisioning.Mode = pb.RemoteProvisioning_MODE_NONE
	devService, err := serviceDevice.New(context.Background(), func() configDevice.Config {
		return cfg.Clients.Device
	}, s.logger)
	if err != nil {
		return err
	}
	s.init(ctx, devService)
	return nil
}

func (s *ClientApplicationServer) Initialize(ctx context.Context, req *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	s.initializationMutex.Lock()
	defer s.initializationMutex.Unlock()

	devService := s.serviceDevice.Load()
	if devService != nil && devService.IsInitialized() {
		return nil, status.Errorf(codes.FailedPrecondition, errAlreadyInitialized)
	}
	if req.GetPreSharedKey() == nil {
		return s.InitializeRemoteProvisioning(ctx, req)
	}
	if req.GetPreSharedKey().GetSubjectId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared subjectId(%v)", req.GetPreSharedKey().GetSubjectId())
	}
	subjectID := events.OwnerToUUID(req.GetPreSharedKey().GetSubjectId())
	_, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared subjectId(%v): %v", req.GetPreSharedKey().GetSubjectId(), err)
	}
	if req.GetPreSharedKey().GetKey() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pre-shared key")
	}
	if err := s.initWithPSK(ctx, req.GetPreSharedKey().GetSubjectId(), req.GetPreSharedKey().GetKey()); err != nil {
		return nil, err
	}
	if _, err := s.ClearCache(ctx, &pb.ClearCacheRequest{}); err != nil {
		log.Warnf("cannot clear device cache: %v", err)
	}
	return &pb.InitializeResponse{}, nil
}

func (s *ClientApplicationServer) FinishInitialize(ctx context.Context, req *pb.FinishInitializeRequest) (*pb.FinishInitializeResponse, error) {
	s.initializationMutex.Lock()
	defer s.initializationMutex.Unlock()

	devService := s.serviceDevice.Load()
	if devService != nil && devService.IsInitialized() {
		return nil, status.Errorf(codes.FailedPrecondition, errAlreadyInitialized)
	}
	if err := s.updateIdentityCertificate(ctx, req); err != nil {
		return nil, err
	}
	if _, err := s.ClearCache(ctx, &pb.ClearCacheRequest{}); err != nil {
		log.Warnf("cannot clear device cache: %v", err)
	}
	return &pb.FinishInitializeResponse{}, nil
}
