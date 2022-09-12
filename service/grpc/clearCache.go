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
	"fmt"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
)

func (s *ClientApplicationServer) ClearCache(ctx context.Context, _ *pb.ClearCacheRequest) (*pb.ClearCacheResponse, error) {
	var errors []error
	s.devices.Range(func(key uuid.UUID, dev *device) bool {
		s.devices.Delete(key)
		err := dev.Close(ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("cannot close device %v connections: %w", key, err))
		}
		return true
	})
	var err error
	switch len(errors) {
	case 0:
	case 1:
		err = errors[0]
	default:
		err = fmt.Errorf("%v", errors)
	}
	if err != nil {
		s.logger.Warnf("cannot properly clear cache: %w", err)
	}

	return &pb.ClearCacheResponse{}, nil
}
