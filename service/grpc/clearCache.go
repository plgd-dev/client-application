package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
)

func (s *DeviceGatewayServer) ClearCache(ctx context.Context, _ *pb.ClearCacheRequest) (*pb.ClearCacheResponse, error) {
	var errors []error
	s.devices.Range(func(key, value interface{}) bool {
		dev, ok := value.(*device)
		if !ok {
			return true
		}
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
