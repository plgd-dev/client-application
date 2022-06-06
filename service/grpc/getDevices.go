package grpc

import (
	"context"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/device/schema"
	"github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	"github.com/plgd-dev/device/schema/resources"
	"github.com/plgd-dev/go-coap/v2/message"
	coapCodes "github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	"github.com/plgd-dev/hub/v2/pkg/log"
	pkgStrings "github.com/plgd-dev/hub/v2/pkg/strings"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Devices []*Device

func (d Devices) Sort() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

type Device struct {
	ID string

	private struct {
		mutex              sync.RWMutex
		ResourceTypes      []string
		Endpoints          schema.Endpoints
		OwnershipStatus    pb.Device_OwnershipStatus
		DeviceResourceBody *pb.Content
	}
}

func (d *Device) ToProto() *pb.Device {
	d.private.mutex.RLock()
	defer d.private.mutex.RUnlock()

	eps := make([]string, 0, len(d.private.Endpoints))
	for _, ep := range d.private.Endpoints {
		eps = append(eps, ep.URI)
	}
	return &pb.Device{
		Content:         d.private.DeviceResourceBody.Clone(),
		OwnershipStatus: d.private.OwnershipStatus,
		Endpoints:       eps,
	}
}

func (d *Device) UpdateDeviceResourceBody(body *pb.Content) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.DeviceResourceBody = body
}

func filterEndpoints(endpoints schema.Endpoints, ipv4TcpEndpoint schema.Endpoint, ipv4UdpEndpoint schema.Endpoint, ipv6TcpEndpoint schema.Endpoint, ipv6UdpEndpoint schema.Endpoint, ipv4secureTcpEndpoint schema.Endpoint, ipv4secureUdpEndpoint schema.Endpoint, ipv6secureTcpEndpoint schema.Endpoint, ipv6secureUdpEndpoint schema.Endpoint) (schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint) {
	for i := range endpoints {
		endpoint := endpoints[i]
		addr, err := endpoint.GetAddr()
		if err != nil {
			continue
		}
		switch addr.GetScheme() {
		case string(schema.TCPScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6TcpEndpoint = endpoint
			} else {
				ipv4TcpEndpoint = endpoint
			}
		case string(schema.UDPScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6UdpEndpoint = endpoint
			} else {
				ipv4UdpEndpoint = endpoint
			}
		case string(schema.TCPSecureScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6secureTcpEndpoint = endpoint
			} else {
				ipv4secureTcpEndpoint = endpoint
			}
		case string(schema.UDPSecureScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6secureUdpEndpoint = endpoint
			} else {
				ipv4secureUdpEndpoint = endpoint
			}
		}
	}
	return ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint
}

func (d *Device) updateEndpointsLocked(endpoints schema.Endpoints) {
	var ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint schema.Endpoint
	ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint = filterEndpoints(d.private.Endpoints, ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint)
	ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint = filterEndpoints(endpoints, ipv4TcpEndpoint, ipv4UdpEndpoint, ipv6TcpEndpoint, ipv6UdpEndpoint, ipv4secureTcpEndpoint, ipv4secureUdpEndpoint, ipv6secureTcpEndpoint, ipv6secureUdpEndpoint)

	newEndpoints := make(schema.Endpoints, 0, 8)
	if _, err := ipv4UdpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4UdpEndpoint)
	}
	if _, err := ipv6UdpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6UdpEndpoint)
	}
	if _, err := ipv4TcpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4TcpEndpoint)
	}
	if _, err := ipv6TcpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6TcpEndpoint)
	}
	if _, err := ipv4secureUdpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4secureUdpEndpoint)
	}
	if _, err := ipv6secureUdpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6secureUdpEndpoint)
	}
	if _, err := ipv4secureTcpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4secureTcpEndpoint)
	}
	if _, err := ipv6secureTcpEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6secureTcpEndpoint)
	}
	d.private.Endpoints = newEndpoints
}

func (d *Device) UpdateDeviceMetadata(resourceTypes []string, endpoints schema.Endpoints, ownershipStatus pb.Device_OwnershipStatus) {
	d.private.mutex.Lock()
	defer d.private.mutex.Unlock()
	d.private.ResourceTypes = resourceTypes
	d.private.OwnershipStatus = ownershipStatus
	d.updateEndpointsLocked(endpoints)
}

func (s *DeviceGatewayServer) discoverMetadataResources(ctx context.Context, discoveryCfg core.DiscoveryConfiguration) error {
	var lock sync.Mutex
	var errors []error

	discoveryClients, err := core.DialDiscoveryAddresses(ctx, discoveryCfg, func(err error) {
		lock.Lock()
		defer lock.Unlock()
		errors = append(errors, err)
	})
	if err != nil {
		return err
	}

	return core.Discover(ctx, discoveryClients, resources.ResourceURI, func(conn *client.ClientConn, req *pool.Message) {
		conn.Close()
		if req.Code() != coapCodes.Content {
			return
		}
		body, err := io.ReadAll(req.Body())
		if err != nil {
			return
		}
		var links schema.ResourceLinks
		if err = cbor.Decode(body, &links); err != nil {
			return
		}
		if len(links) == 0 {
			return
		}
		var deviceID string
		var endpoints schema.Endpoints
		var resourceTypes []string
		ownershipStatus := pb.Device_UNSUPPORTED
		for _, l := range links {
			deviceID = l.GetDeviceID()
			if pkgStrings.Contains(l.ResourceTypes, device.ResourceType) {
				endpoints = l.Endpoints
				resourceTypes = l.ResourceTypes
			}
			if pkgStrings.Contains(l.ResourceTypes, doxm.ResourceType) {
				if len(l.Endpoints.FilterUnsecureEndpoints()) == 0 {
					ownershipStatus = pb.Device_OWNED
				} else {
					ownershipStatus = pb.Device_UNOWNED
				}
			}
		}
		device := Device{
			ID: deviceID,
		}
		device.private.ResourceTypes = resourceTypes
		device.private.Endpoints = endpoints
		device.private.OwnershipStatus = ownershipStatus
		v, loaded := s.devices.LoadOrStore(deviceID, &device)
		if !loaded {
			return
		}
		v.(*Device).UpdateDeviceMetadata(resourceTypes, endpoints, ownershipStatus)
	}, coap.WithResourceType(device.ResourceType), coap.WithResourceType(doxm.ResourceType))
}

func (s *DeviceGatewayServer) discoverDeviceResource(ctx context.Context, discoveryCfg core.DiscoveryConfiguration) error {
	discoveryClients, err := core.DialDiscoveryAddresses(ctx, discoveryCfg, func(err error) {})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return core.Discover(ctx, discoveryClients, device.ResourceURI, func(conn *client.ClientConn, req *pool.Message) {
		conn.Close()
		if req.Code() != coapCodes.Content {
			return
		}
		body, err := io.ReadAll(req.Body())
		if err != nil {
			return
		}
		var d device.Device
		if err = cbor.Decode(body, &d); err != nil {
			return
		}
		if d.ID == "" {
			return
		}
		device := Device{
			ID: d.ID,
		}

		contentFormat, err := req.ContentFormat()
		if err != nil {
			contentFormat = message.AppOcfCbor
		}
		device.private.ResourceTypes = d.ResourceTypes
		device.private.DeviceResourceBody = &pb.Content{
			ContentType: contentFormat.String(),
			Data:        body,
		}
		v, loaded := s.devices.LoadOrStore(d.ID, &device)
		if !loaded {
			return
		}
		v.(*Device).UpdateDeviceResourceBody(device.private.DeviceResourceBody)
	})
}

func (s *DeviceGatewayServer) GetDevices(req *pb.GetDevicesRequest, srv pb.DeviceGateway_GetDevicesServer) error {
	ctx := srv.Context()
	discoveryCfg := core.DefaultDiscoveryConfiguration()
	discoveryCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	errChan := make(chan error, 1)
	go func() {
		defer wg.Done()
		err := s.discoverDeviceResource(discoveryCtx, discoveryCfg)
		if err != nil {
			errChan <- err
		}
	}()
	err := s.discoverMetadataResources(discoveryCtx, discoveryCfg)
	if err != nil {
		log.Errorf("failed to discover device resources: %w", err)
	}
	wg.Wait()
	select {
	case err = <-errChan:
		log.Errorf("failed to discover endpoints and ownership status: %w", err)
	default:
	}

	devices := make(Devices, 0, 128)
	s.devices.Range(func(key, value any) bool {
		devices = append(devices, value.(*Device))
		return true
	})
	devices.Sort()
	for _, device := range devices {
		d := device.ToProto()
		if d.GetContent() == nil {
			continue
		}
		if srv.Send(d) != nil {
			return err
		}
	}

	return nil
}
