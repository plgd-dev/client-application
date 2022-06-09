package grpc

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/device/schema"
	plgdDevice "github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	"github.com/plgd-dev/device/schema/resources"
	"github.com/plgd-dev/go-coap/v2/message"
	coapCodes "github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/log"
	pkgStrings "github.com/plgd-dev/hub/v2/pkg/strings"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	pkgNet "github.com/plgd-dev/kit/v2/net"
	kitStrings "github.com/plgd-dev/kit/v2/strings"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultTimeout = 1000
	MulticastPort  = 5683
)

func filterEndpoints(endpoints schema.Endpoints, ipv4TCPEndpoint schema.Endpoint, ipv4UDPEndpoint schema.Endpoint, ipv6TCPEndpoint schema.Endpoint, ipv6UDPEndpoint schema.Endpoint, ipv4secureTCPEndpoint schema.Endpoint, ipv4secureUDPEndpoint schema.Endpoint, ipv6secureTCPEndpoint schema.Endpoint, ipv6secureUDPEndpoint schema.Endpoint) (schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint, schema.Endpoint) {
	for i := range endpoints {
		endpoint := endpoints[i]
		addr, err := endpoint.GetAddr()
		if err != nil {
			continue
		}
		switch addr.GetScheme() {
		case string(schema.TCPScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6TCPEndpoint = endpoint
			} else {
				ipv4TCPEndpoint = endpoint
			}
		case string(schema.UDPScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6UDPEndpoint = endpoint
			} else {
				ipv4UDPEndpoint = endpoint
			}
		case string(schema.TCPSecureScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6secureTCPEndpoint = endpoint
			} else {
				ipv4secureTCPEndpoint = endpoint
			}
		case string(schema.UDPSecureScheme):
			if strings.ContainsAny(addr.GetHostname(), ":") {
				ipv6secureUDPEndpoint = endpoint
			} else {
				ipv4secureUDPEndpoint = endpoint
			}
		}
	}

	return ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint
}

func (d *device) updateEndpointsLocked(endpoints schema.Endpoints) {
	var ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint schema.Endpoint
	ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint = filterEndpoints(d.private.Endpoints, ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint)
	ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint = filterEndpoints(endpoints, ipv4TCPEndpoint, ipv4UDPEndpoint, ipv6TCPEndpoint, ipv6UDPEndpoint, ipv4secureTCPEndpoint, ipv4secureUDPEndpoint, ipv6secureTCPEndpoint, ipv6secureUDPEndpoint)

	newEndpoints := make(schema.Endpoints, 0, 8)
	if _, err := ipv4UDPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4UDPEndpoint)
	}
	if _, err := ipv6UDPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6UDPEndpoint)
	}
	if _, err := ipv4TCPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4TCPEndpoint)
	}
	if _, err := ipv6TCPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6TCPEndpoint)
	}
	if _, err := ipv4secureUDPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4secureUDPEndpoint)
	}
	if _, err := ipv6secureUDPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6secureUDPEndpoint)
	}
	if _, err := ipv4secureTCPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv4secureTCPEndpoint)
	}
	if _, err := ipv6secureTCPEndpoint.GetAddr(); err == nil {
		newEndpoints = append(newEndpoints, ipv6secureTCPEndpoint)
	}
	d.private.Endpoints = newEndpoints
}

func processDiscoveryResourceResponse(resp *pool.Message) (*device, error) {
	if resp.Code() != coapCodes.Content {
		return nil, fmt.Errorf("unexpected response code: %d", resp.Code())
	}
	body, err := io.ReadAll(resp.Body())
	if err != nil {
		return nil, err
	}
	var links schema.ResourceLinks
	if err = cbor.Decode(body, &links); err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("no links in response")
	}
	var deviceID string
	var endpoints schema.Endpoints
	var resourceTypes []string
	ownershipStatus := grpcgwPb.Device_UNSUPPORTED
	for _, l := range links {
		deviceID = l.GetDeviceID()
		if pkgStrings.Contains(l.ResourceTypes, plgdDevice.ResourceType) {
			endpoints = l.Endpoints
			resourceTypes = l.ResourceTypes
		}
		if pkgStrings.Contains(l.ResourceTypes, doxm.ResourceType) {
			if len(l.Endpoints.FilterUnsecureEndpoints()) == 0 {
				ownershipStatus = grpcgwPb.Device_OWNED
			} else {
				ownershipStatus = grpcgwPb.Device_UNOWNED
			}
		}
	}
	device := newDevice(deviceID)
	device.private.ResourceTypes = resourceTypes
	device.updateEndpointsLocked(endpoints)
	device.private.OwnershipStatus = ownershipStatus

	return device, nil
}

func onDiscoveryResourceResponse(resp *pool.Message, devices *sync.Map) error {
	dev, err := processDiscoveryResourceResponse(resp)
	if err != nil {
		return err
	}
	v, loaded := devices.LoadOrStore(dev.ID, dev)
	if !loaded {
		return nil
	}
	if d, ok := v.(*device); ok {
		d.UpdateDeviceMetadata(dev.private.ResourceTypes, dev.private.Endpoints, dev.private.OwnershipStatus)
	}

	return nil
}

func discoverDiscoveryResources(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onResponse func(conn *client.ClientConn, resp *pool.Message)) error {
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
	defer func() {
		for _, c := range discoveryClients {
			_ = c.Close()
		}
	}()
	if len(errors) > 0 {
		lock.Lock()
		dbgErrors := errors
		errors = nil
		lock.Unlock()
		log.Debugf("some fails occurs during discover discovery resources of the device: %v", dbgErrors)
	}

	return core.Discover(ctx, discoveryClients, resources.ResourceURI, onResponse, coap.WithResourceType(plgdDevice.ResourceType), coap.WithResourceType(doxm.ResourceType))
}

func processDeviceResourceResponse(resp *pool.Message) (*device, error) {
	if resp.Code() != coapCodes.Content {
		return nil, fmt.Errorf("unexpected response code: %s", resp.Code())
	}
	body, err := io.ReadAll(resp.Body())
	if err != nil {
		return nil, err
	}
	var d plgdDevice.Device
	if err = cbor.Decode(body, &d); err != nil {
		return nil, err
	}
	if d.ID == "" {
		return nil, fmt.Errorf("device ID is empty")
	}
	device := newDevice(d.ID)
	contentFormat, err := resp.ContentFormat()
	if err != nil {
		contentFormat = message.AppOcfCbor
	}
	device.private.ResourceTypes = d.ResourceTypes
	device.private.DeviceResourceBody = &commands.Content{
		ContentType: contentFormat.String(),
		Data:        body,
	}

	return device, nil
}

func onDeviceResourceResponse(resp *pool.Message, devices *sync.Map) error {
	dev, err := processDeviceResourceResponse(resp)
	if err != nil {
		return err
	}
	v, loaded := devices.LoadOrStore(dev.ID, dev)
	if !loaded {
		return nil
	}
	if d, ok := v.(*device); ok {
		d.UpdateDeviceResourceBody(dev.private.DeviceResourceBody)
	}

	return nil
}

func discoverDeviceResource(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onResponse func(conn *client.ClientConn, resp *pool.Message)) error {
	discoveryClients, err := core.DialDiscoveryAddresses(ctx, discoveryCfg, func(err error) {})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer func() {
		for _, c := range discoveryClients {
			_ = c.Close()
		}
	}()

	return core.Discover(ctx, discoveryClients, plgdDevice.ResourceURI, onResponse)
}

func toDiscoveryConfiguration(ipVersionFilter ipVersionFilter) core.DiscoveryConfiguration {
	discoveryCfg := core.DefaultDiscoveryConfiguration()
	if ipVersionFilter&ipv4 == 0 {
		discoveryCfg.MulticastAddressUDP4 = nil
	}
	if ipVersionFilter&ipv6 == 0 {
		discoveryCfg.MulticastAddressUDP6 = nil
	}
	return discoveryCfg
}

func getDevicesByMulticast(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onDeviceResourceResponse, onDiscoveryResourceResponse func(conn *client.ClientConn, resp *pool.Message)) {
	var wg sync.WaitGroup
	wg.Add(1)

	errChan := make(chan error, 1)
	go func() {
		defer wg.Done()
		err := discoverDeviceResource(ctx, discoveryCfg, onDeviceResourceResponse)
		if err != nil {
			errChan <- err
		}
	}()
	err := discoverDiscoveryResources(ctx, discoveryCfg, onDiscoveryResourceResponse)
	if err != nil {
		log.Errorf("failed to discover device resources: %w", err)
	}
	wg.Wait()
	select {
	case err = <-errChan:
		log.Errorf("failed to discover endpoints and ownership status: %w", err)
	default:
	}
}

func normalizeEndpoint(endpoint string) (pkgNet.Addr, error) {
	addressPort := endpoint
	addr, err := pkgNet.ParseString(string(schema.UDPScheme), addressPort)
	if err != nil && strings.Contains(err.Error(), "missing port in address") {
		addr, err = pkgNet.ParseString(string(schema.UDPScheme), fmt.Sprintf("%v:%v", addressPort, MulticastPort))
	}
	if err != nil {
		return pkgNet.Addr{}, fmt.Errorf("invalid endpoint: %s", endpoint)
	}
	return addr, nil
}

func getDeviceByMulticastAddress(ctx context.Context, addr pkgNet.Addr, devices *sync.Map) error {
	address := fmt.Sprintf("%s:%d", addr.GetHostname(), addr.GetPort())
	multicastAddr := []string{address}
	discoveryConfiguration := core.DefaultDiscoveryConfiguration()
	if strings.Contains(addr.GetHostname(), ":") {
		discoveryConfiguration.MulticastAddressUDP6 = multicastAddr
		discoveryConfiguration.MulticastAddressUDP4 = nil
	} else {
		discoveryConfiguration.MulticastAddressUDP4 = multicastAddr
		discoveryConfiguration.MulticastAddressUDP6 = nil
	}
	deviceResource := atomic.NewBool(false)
	discoveryResource := atomic.NewBool(false)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	getDevicesByMulticast(ctx, discoveryConfiguration, func(conn *client.ClientConn, resp *pool.Message) {
		_ = conn.Close()
		err := onDeviceResourceResponse(resp, devices)
		if err != nil {
			return
		}
		if deviceResource.CAS(false, true) && discoveryResource.Load() {
			cancel()
		}
	}, func(conn *client.ClientConn, resp *pool.Message) {
		_ = conn.Close()
		err := onDiscoveryResourceResponse(resp, devices)
		if err != nil {
			return
		}
		if discoveryResource.CAS(false, true) && deviceResource.Load() {
			cancel()
		}
	})
	return nil
}

func getDeviceByAddress(ctx context.Context, addr pkgNet.Addr, devices *sync.Map) error {
	if addr.GetPort() == MulticastPort {
		return getDeviceByMulticastAddress(ctx, addr, devices)
	}
	address := fmt.Sprintf("%s:%d", addr.GetHostname(), addr.GetPort())
	client, err := udp.Dial(address, udp.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()
	resp, err := client.Get(ctx, plgdDevice.ResourceURI)
	if err != nil {
		return err
	}
	deviceRes, err := processDeviceResourceResponse(resp)
	if err != nil {
		return err
	}
	opts := make(message.Options, 0, 2)
	coap.WithResourceType(plgdDevice.ResourceType)(opts)
	coap.WithResourceType(doxm.ResourceType)(opts)
	resp, err = client.Get(ctx, resources.ResourceURI, opts...)
	if err != nil {
		return err
	}
	discoveryRes, err := processDiscoveryResourceResponse(resp)
	if err != nil {
		return err
	}
	v, loaded := devices.LoadOrStore(deviceRes.ID, deviceRes)
	if d, ok := v.(*device); ok {
		if loaded {
			d.UpdateDeviceResourceBody(deviceRes.private.DeviceResourceBody)
		}
		d.UpdateDeviceMetadata(discoveryRes.private.ResourceTypes, discoveryRes.private.Endpoints, discoveryRes.private.OwnershipStatus)
	}
	return nil
}

func getDevicesByEndpoints(ctx context.Context, endpoints []string, devices *sync.Map) {
	addresses := make([]pkgNet.Addr, 0, len(endpoints))
	for _, endpoint := range endpoints {
		addr, err := normalizeEndpoint(endpoint)
		if err != nil {
			log.Errorf("%w", err)

			continue
		}
		addresses = append(addresses, addr)
	}
	if len(addresses) == 1 {
		err := getDeviceByAddress(ctx, addresses[0], devices)
		if err != nil {
			log.Errorf("cannot get device by address %v: %w", addresses[0], err)
		}
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(endpoints))
	for _, addr := range addresses {
		go func(addr pkgNet.Addr) {
			defer wg.Done()
			err := getDeviceByAddress(ctx, addr, devices)
			if err != nil {
				log.Errorf("cannot get device by address %v: %w", addr, err)
			}
		}(addr)
	}
	wg.Wait()
}

type ipVersionFilter int8

const (
	ipv4 ipVersionFilter = 1 << iota
	ipv6
)

func toUseMulticastFilter(v []pb.GetDevicesRequest_UseMulticast) ipVersionFilter {
	var f ipVersionFilter
	for _, ipv := range v {
		switch ipv {
		case pb.GetDevicesRequest_IPV4:
			f |= ipv4
		case pb.GetDevicesRequest_IPV6:
			f |= ipv6
		}
	}
	return f
}

func filterByType(device *grpcgwPb.Device, filteredTypes []string) bool {
	if len(filteredTypes) == 0 {
		return true
	}
	types := kitStrings.MakeSet(filteredTypes...)
	return types.HasOneOf(device.GetTypes()...)
}

func filterByOwnershipStatus(device *grpcgwPb.Device, filteredOwnershipStatus []pb.GetDevicesRequest_OwnershipStatusFilter) bool {
	if len(filteredOwnershipStatus) == 0 {
		return true
	}
	for _, status := range filteredOwnershipStatus {
		switch status {
		case pb.GetDevicesRequest_OWNED:
			if device.GetOwnershipStatus() == grpcgwPb.Device_OWNED {
				return true
			}
		case pb.GetDevicesRequest_UNOWNED:
			if device.GetOwnershipStatus() == grpcgwPb.Device_UNOWNED {
				return true
			}
		}
	}
	return false
}

// If use_cache, use_multicast, use_endpoints are not set, then it will set use_multicast with [IPV4,IPV6].
func tryToSetDefaultRequest(req *pb.GetDevicesRequest) *pb.GetDevicesRequest {
	if req == nil {
		req = &pb.GetDevicesRequest{}
	}
	if !req.GetUseCache() && len(req.GetUseMulticast()) == 0 && len(req.GetUseEndpoints()) == 0 {
		req.UseMulticast = []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4, pb.GetDevicesRequest_IPV6}
	}
	return req
}

func (s *DeviceGatewayServer) GetDevices(req *pb.GetDevicesRequest, srv pb.DeviceGateway_GetDevicesServer) error {
	req = tryToSetDefaultRequest(req)
	ctx := srv.Context()
	var toCall []func()
	var discoveredDevices sync.Map
	var cachedDevices sync.Map
	if req.GetTimeout() == 0 {
		req.Timeout = DefaultTimeout
	}
	discoveryCtx, cancel := context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Millisecond)
	defer cancel()
	if req.UseCache {
		s.devices.Range(func(key, value interface{}) bool {
			cachedDevices.Store(key, value)
			return true
		})
	}
	if len(req.GetUseMulticast()) > 0 {
		toCall = append(toCall, func() {
			getDevicesByMulticast(discoveryCtx, toDiscoveryConfiguration(toUseMulticastFilter(req.GetUseMulticast())), func(conn *client.ClientConn, resp *pool.Message) {
				_ = conn.Close()
				_ = onDeviceResourceResponse(resp, &discoveredDevices)
			}, func(conn *client.ClientConn, resp *pool.Message) {
				_ = conn.Close()
				_ = onDiscoveryResourceResponse(resp, &discoveredDevices)
			})
		},
		)
	}
	if len(req.GetUseEndpoints()) > 0 {
		toCall = append(toCall, func() {
			getDevicesByEndpoints(discoveryCtx, req.GetUseEndpoints(), &discoveredDevices)
		})
	}

	var wg sync.WaitGroup
	wg.Add(len(toCall))
	for _, f := range toCall {
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}
	wg.Wait()

	devs := make(devices, 0, 128)
	discoveredDevices.Range(func(key, value any) bool {
		if d, ok := value.(*device); ok {
			devs = append(devs, d)
		}
		s.devices.Store(key, value)
		cachedDevices.Delete(key)
		return true
	})
	cachedDevices.Range(func(key, value any) bool {
		if d, ok := value.(*device); ok {
			devs = append(devs, d)
		}
		return true
	})
	devs.Sort()
	for _, device := range devs {
		d := device.ToProto()
		if d.GetData().GetContent() == nil {
			continue
		}
		if !filterByType(d, req.GetTypeFilter()) {
			continue
		}
		if !filterByOwnershipStatus(d, req.GetOwnershipStatusFilter()) {
			continue
		}
		if err := srv.Send(d); err != nil {
			return err
		}
	}

	return nil
}
