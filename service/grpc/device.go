package grpc

import (
	"context"
	"fmt"
	"sort"
	"sync"

	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/schema"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type devices []*device

func (d devices) Sort() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

type device struct {
	ID     string
	logger log.Logger

	private struct {
		mutex              sync.RWMutex
		ResourceTypes      []string
		Endpoints          schema.Endpoints
		OwnershipStatus    grpcgwPb.Device_OwnershipStatus
		DeviceResourceBody *commands.Content
		api                *core.Device
	}
	*core.Device
}

func newDevice(deviceID string, serviceDevice *serviceDevice.Service, logger log.Logger) *device {
	d := device{
		ID:     deviceID,
		logger: logger.With(log.DeviceIDKey, deviceID),
	}
	coreDeviceCfg := serviceDevice.GetDeviceConfiguration()
	coreDeviceCfg.ErrFunc = d.ErrorFunc
	d.Device = core.NewDevice(coreDeviceCfg, d.ID, []string{}, d.GetEndpoints)
	return &d
}

func (d *device) ErrorFunc(err error) {
	d.logger.Debug(err)
}

func (d *device) ToProto() *grpcgwPb.Device {
	d.private.mutex.RLock()
	defer d.private.mutex.RUnlock()

	eps := make([]string, 0, len(d.private.Endpoints))
	for _, ep := range d.private.Endpoints {
		eps = append(eps, ep.URI)
	}

	return &grpcgwPb.Device{
		Id:    d.ID,
		Types: d.private.ResourceTypes,
		Data: &events.ResourceChanged{
			Content: &commands.Content{
				Data:              d.private.DeviceResourceBody.GetData(),
				ContentType:       d.private.DeviceResourceBody.GetContentType(),
				CoapContentFormat: d.private.DeviceResourceBody.GetCoapContentFormat(),
			},
			Status: commands.Status_OK,
		},
		OwnershipStatus: d.private.OwnershipStatus,
		Endpoints:       eps,
	}
}

func (d *device) GetEndpoints() schema.Endpoints {
	d.private.mutex.RLock()
	defer d.private.mutex.RUnlock()
	endpoints := make(schema.Endpoints, len(d.private.Endpoints))
	copy(endpoints, d.private.Endpoints)
	return endpoints
}

func normalizeHref(href string) string {
	if href == "" {
		return ""
	}
	if href[0] == '/' {
		return href
	}
	return "/" + href
}

func getOwnershipStatus(links schema.ResourceLinks) grpcgwPb.Device_OwnershipStatus {
	doxmRes := links.GetResourceLinks(doxm.ResourceType)
	if len(doxmRes) == 0 {
		return grpcgwPb.Device_UNSUPPORTED
	}
	l := doxmRes[0]
	if len(l.Endpoints.FilterUnsecureEndpoints()) == 0 {
		return grpcgwPb.Device_OWNED
	}
	return grpcgwPb.Device_UNOWNED
}

func (d *device) getResourceLinksAndRefreshCache(ctx context.Context) (schema.ResourceLinks, error) {
	if d.Device == nil {
		return nil, status.Error(codes.Internal, "device is not initialized")
	}
	links, err := d.GetResourceLinks(ctx, d.GetEndpoints())
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource links for device %v: %w", d.ID, err)).Err()
	}
	d.updateOwnershipStatus(getOwnershipStatus(links))
	devLinks := links.GetResourceLinks(plgdDevice.ResourceType)
	if len(devLinks) > 0 {
		d.updateResourceTypes(devLinks[0].ResourceTypes)
	}
	return links, nil
}

func (d *device) getResourceLink(ctx context.Context, resourceID *commands.ResourceId) (schema.ResourceLink, error) {
	links, err := d.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return schema.ResourceLink{}, nil
	}
	link, ok := links.GetResourceLink(normalizeHref(resourceID.GetHref()))
	if !ok {
		return schema.ResourceLink{}, status.Errorf(codes.NotFound, "cannot find resource link %v for device %v", resourceID.GetHref(), d.ID)
	}
	return link, nil
}

func (d *device) updateDeviceMetadata(resourceTypes []string, endpoints schema.Endpoints, ownershipStatus grpcgwPb.Device_OwnershipStatus) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.ResourceTypes = resourceTypes
	d.private.OwnershipStatus = ownershipStatus
	d.updateEndpointsLocked(endpoints)
}

func (d *device) updateResourceTypes(resourceTypes []string) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.ResourceTypes = resourceTypes
}

func (d *device) updateOwnershipStatus(ownershipStatus grpcgwPb.Device_OwnershipStatus) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.OwnershipStatus = ownershipStatus
}

func (d *device) updateDeviceResourceBody(body *commands.Content) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.DeviceResourceBody = body
}