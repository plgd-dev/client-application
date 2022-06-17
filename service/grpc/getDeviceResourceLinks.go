package grpc

import (
	"context"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
)

func (s *ClientApplicationServer) GetDeviceResourceLinks(ctx context.Context, req *pb.GetDeviceResourceLinksRequest) (*events.ResourceLinksPublished, error) {
	dev, err := s.getDevice(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	return &events.ResourceLinksPublished{
		Resources: commands.SchemaResourceLinksToResources(links, time.Time{}),
		DeviceId:  dev.DeviceID(),
	}, nil
}
