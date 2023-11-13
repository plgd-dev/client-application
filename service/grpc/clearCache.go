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
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/plgd-dev/client-application/pb"
)

func closeDevice(dev *device) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	return dev.Close(ctx)
}

func closeDevices(devices map[uuid.UUID]*device) error {
	var errors *multierror.Error
	for key, dev := range devices {
		err := closeDevice(dev)
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("cannot close device %v connections: %w", key, err))
		}
	}
	if errors == nil {
		return nil
	}
	switch errors.Len() {
	case 0:
		return nil
	case 1:
		return errors.Errors[0]
	default:
		return errors
	}
}

func (s *ClientApplicationServer) ClearCache(ctx context.Context, _ *pb.ClearCacheRequest) (*pb.ClearCacheResponse, error) {
	devices := s.devices.LoadAndDeleteAll()
	go func(devices map[uuid.UUID]*device) {
		err := closeDevices(devices)
		if err != nil {
			s.logger.Warnf("cannot properly clear cache: %w", err)
		}
	}(devices)
	return &pb.ClearCacheResponse{}, nil
}
