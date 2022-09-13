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
	"github.com/plgd-dev/hub/v2/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ClientApplicationServer) reset(ctx context.Context) {
	s.jwksCache.Store(nil)
	s.serviceDevice.Reset()
	s.csrCache.DeleteAll()
	s.remoteOwnSignCache.Range(func(key uuid.UUID, value *remoteSign) bool {
		s.remoteOwnSignCache.Delete(key)
		value.cancel()
		return true
	})
	_, err := s.ClearCache(ctx, &pb.ClearCacheRequest{})
	if err != nil {
		log.Warnf("cannot clear cache: %v", err)
	}
}

func (s *ClientApplicationServer) Reset(ctx context.Context, req *pb.ResetRequest) (*pb.ResetResponse, error) {
	if !s.serviceDevice.IsInitialized() {
		return nil, status.Errorf(codes.FailedPrecondition, "not initialized")
	}
	s.reset(ctx)
	return &pb.ResetResponse{}, nil
}
