package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/device/schema"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"github.com/plgd-dev/kit/v2/log"
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
	ID string

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

func dialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithKeepAlive(time.Minute),
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, append(dialOpts, opts...)...)
}

func dialUDP(ctx context.Context, addr string, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithKeepAlive(time.Minute),
	}
	return coap.DialUDP(ctx, addr, append(dialOpts, opts...)...)
}

func dialTCP(ctx context.Context, addr string, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithKeepAlive(time.Minute),
	}
	return coap.DialTCP(ctx, addr, append(dialOpts, opts...)...)
}

func dialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithKeepAlive(time.Minute),
	}
	return coap.DialTCPSecure(ctx, addr, tlsCfg, append(dialOpts, opts...)...)
}

func errFunc(err error) {
	log.Debug(err)
}

func newDevice(deviceID string) *device {
	d := device{
		ID: deviceID,
	}

	d.Device = core.NewDevice(core.DeviceConfiguration{
		DialDTLS: dialDTLS,
		DialUDP:  dialUDP,
		DialTCP:  dialTCP,
		DialTLS:  dialTLS,
		ErrFunc:  errFunc,
		TLSConfig: &core.TLSConfig{
			GetCertificate: func() (tls.Certificate, error) {
				return tls.Certificate{}, fmt.Errorf("not implemented")
			},
			GetCertificateAuthorities: func() ([]*x509.Certificate, error) {
				return nil, fmt.Errorf("not implemented")
			},
		},
	}, d.ID, []string{}, d.GetEndpoints)
	return &d
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

func (d *device) getResourceLink(ctx context.Context, resourceID *commands.ResourceId) (schema.ResourceLink, error) {
	if d.Device == nil {
		return schema.ResourceLink{}, status.Error(codes.Internal, "device is not initialized")
	}
	links, err := d.GetResourceLinks(ctx, d.GetEndpoints())
	if err != nil {
		return schema.ResourceLink{}, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource links for device %v: %w", d.ID, err)).Err()
	}
	link, ok := links.GetResourceLink(normalizeHref(resourceID.GetHref()))
	if !ok {
		return schema.ResourceLink{}, status.Errorf(codes.NotFound, "cannot find resource link %v for device %v", resourceID.GetHref(), d.ID)
	}
	return link, nil
}

func (d *device) UpdateDeviceMetadata(resourceTypes []string, endpoints schema.Endpoints, ownershipStatus grpcgwPb.Device_OwnershipStatus) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.ResourceTypes = resourceTypes
	d.private.OwnershipStatus = ownershipStatus
	d.updateEndpointsLocked(endpoints)
}

func (d *device) UpdateDeviceResourceBody(body *commands.Content) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.DeviceResourceBody = body
}
