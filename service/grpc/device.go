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
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/device/v2/schema"
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	"github.com/plgd-dev/device/v2/schema/doxm"
	"github.com/plgd-dev/device/v2/schema/interfaces"
	"github.com/plgd-dev/device/v2/schema/resources"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/pkg/strings"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type devices []*device

func (d devices) Sort() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID.String() < d[j].ID.String()
	})
}

type device struct {
	ID     uuid.UUID
	logger log.Logger

	private struct {
		mutex              sync.RWMutex
		ResourceTypes      []string
		DeviceURI          string
		Endpoints          schema.Endpoints
		OwnershipStatus    grpcgwPb.Device_OwnershipStatus
		DeviceResourceBody *commands.Content
		api                *core.Device
	}
	*core.Device
}

func newDevice(deviceID uuid.UUID, serviceDevice *serviceDevice.Service, logger log.Logger) *device {
	coreDeviceCfg := serviceDevice.GetDeviceConfiguration()
	d := device{
		ID:     deviceID,
		logger: logger.With(log.DeviceIDKey, deviceID),
	}
	coreDeviceCfg.Logger = serviceDevice.DeviceLogger()
	d.Device = core.NewDevice(coreDeviceCfg, deviceID.String(), []string{}, d.GetEndpoints)
	return &d
}

func (d *device) ErrorFunc(err error) {
	d.logger.Debug(err)
}

func (d *device) hasDeviceResourceBody() bool {
	d.private.mutex.RLock()
	defer d.private.mutex.RUnlock()
	return d.private.DeviceResourceBody != nil
}

func (d *device) ToProto() *grpcgwPb.Device {
	d.private.mutex.RLock()
	defer d.private.mutex.RUnlock()

	eps := make([]string, 0, len(d.private.Endpoints))
	for _, ep := range d.private.Endpoints {
		eps = append(eps, ep.URI)
	}

	return &grpcgwPb.Device{
		Id:    d.ID.String(),
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

func getOwnershipStatus(l schema.ResourceLink) grpcgwPb.Device_OwnershipStatus {
	if len(l.Endpoints.FilterUnsecureEndpoints()) == 0 {
		return grpcgwPb.Device_OWNED
	}
	return grpcgwPb.Device_UNOWNED
}

func getOwnershipStatusLinks(links schema.ResourceLinks) grpcgwPb.Device_OwnershipStatus {
	doxmRes := links.GetResourceLinks(doxm.ResourceType)
	if len(doxmRes) == 0 {
		return grpcgwPb.Device_UNSUPPORTED
	}
	return getOwnershipStatus(doxmRes[0])
}

func (d *device) getResourceLinksAndRefreshCache(ctx context.Context) (schema.ResourceLinks, error) {
	if d.Device == nil {
		return nil, status.Error(codes.Internal, "device is not initialized")
	}
	links, err := d.GetResourceLinks(ctx, d.GetEndpoints(), coap.WithDeviceID(d.DeviceID()))
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource links for device %v: %w", d.ID, err)).Err()
	}
	d.updateOwnershipStatus(getOwnershipStatusLinks(links))
	devLinks := links.GetResourceLinks(plgdDevice.ResourceType)
	if len(devLinks) > 0 {
		d.updateResourceTypes(devLinks[0].ResourceTypes)
	}
	return links, nil
}

func (d *device) getResourceLink(ctx context.Context, resourceID *commands.ResourceId) (schema.ResourceLink, error) {
	links, err := d.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return schema.ResourceLink{}, err
	}
	link, ok := links.GetResourceLink(normalizeHref(resourceID.GetHref()))
	if !ok {
		return schema.ResourceLink{}, status.Errorf(codes.NotFound, "cannot find resource link %v for device %v", resourceID.GetHref(), d.ID)
	}
	return link, nil
}

func (d *device) checkAccess(link schema.ResourceLink) error {
	if d.ToProto().GetOwnershipStatus() != grpcgwPb.Device_OWNED && len(link.Endpoints.FilterUnsecureEndpoints()) == 0 {
		return status.Error(codes.PermissionDenied, "device is not owned")
	}
	return nil
}

func (d *device) getResourceLinkAndCheckAccess(ctx context.Context, resourceID *commands.ResourceId, resInterface string) (schema.ResourceLink, error) {
	link, err := d.getResourceLink(ctx, resourceID)
	if err != nil {
		return link, err
	}
	if strings.Contains(link.ResourceTypes, resources.ResourceType) && resInterface == interfaces.OC_IF_B {
		link.Endpoints = link.Endpoints.FilterSecureEndpoints()
	}
	return link, d.checkAccess(link)
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

func (d *device) update(data *device) {
	data.private.mutex.RLock()
	defer data.private.mutex.RUnlock()
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.DeviceResourceBody = data.private.DeviceResourceBody
	d.private.ResourceTypes = data.private.ResourceTypes
	d.private.OwnershipStatus = data.private.OwnershipStatus
	d.updateEndpointsLocked(data.private.Endpoints)
}

func (d *device) provision(ctx context.Context, links schema.ResourceLinks, action func(context.Context, *core.ProvisioningClient) error) (err error) {
	var p *core.ProvisioningClient
	p, err = d.Provision(ctx, links)
	if err != nil {
		return err
	}
	defer func() {
		deadline, ok := ctx.Deadline()
		if ctx.Err() != nil || !ok || time.Until(deadline) < time.Second {
			ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			ctx = ctx1
		}
		pErr := p.Close(ctx)
		if err == nil {
			err = pErr
		}
	}()
	err = action(ctx, p)
	return
}
