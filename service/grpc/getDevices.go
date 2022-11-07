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
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/device/v2/schema"
	plgdDevice "github.com/plgd-dev/device/v2/schema/device"
	"github.com/plgd-dev/device/v2/schema/doxm"
	"github.com/plgd-dev/device/v2/schema/resources"
	"github.com/plgd-dev/go-coap/v3/message"
	coapCodes "github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/options"
	coapSync "github.com/plgd-dev/go-coap/v3/pkg/sync"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/go-coap/v3/udp/client"
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
	DefaultTimeout = 2 * time.Second
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

func processDiscoveryResourceResponse(serviceDevice *serviceDevice.Service, logger log.Logger, remoteAddr net.Addr, resp *pool.Message) (*device, error) {
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

	addr, err := pkgNet.ParseString(string(schema.UDPScheme), remoteAddr.String())
	if err != nil {
		return nil, err
	}
	links = links.PatchEndpoint(addr, schema.Endpoints{})

	var deviceID string
	var endpoints schema.Endpoints
	var resourceTypes []string
	ownershipStatus := getOwnershipStatus(links)
	for _, l := range links {
		deviceID = l.GetDeviceID()
		if pkgStrings.Contains(l.ResourceTypes, plgdDevice.ResourceType) {
			endpoints = l.Endpoints
			resourceTypes = l.ResourceTypes
		}
	}
	devID, err := uuid.Parse(deviceID)
	if err != nil {
		return nil, fmt.Errorf("cannot parse device ID('%v'): %w", devID, err)
	}
	device := newDevice(devID, serviceDevice, logger)
	device.private.ResourceTypes = resourceTypes
	device.updateEndpointsLocked(endpoints)
	device.private.OwnershipStatus = ownershipStatus

	return device, nil
}

func onDiscoveryResourceResponse(serviceDevice *serviceDevice.Service, logger log.Logger, remoteAddr net.Addr, resp *pool.Message, devices *coapSync.Map[uuid.UUID, *device]) error {
	dev, err := processDiscoveryResourceResponse(serviceDevice, logger, remoteAddr, resp)
	if err != nil {
		return err
	}
	d, loaded := devices.LoadOrStore(dev.ID, dev)
	if !loaded {
		return nil
	}
	d.updateDeviceMetadata(dev.private.ResourceTypes, dev.private.Endpoints, dev.private.OwnershipStatus)

	return nil
}

func discoverDiscoveryResources(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onResponse func(conn *client.Conn, resp *pool.Message)) error {
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

func processDeviceResourceResponse(serviceDevice *serviceDevice.Service, logger log.Logger, resp *pool.Message) (*device, error) {
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
	devID, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot parse device ID('%v'): %w", d.ID, err)
	}
	device := newDevice(devID, serviceDevice, logger)
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

func onDeviceResourceResponse(serviceDevice *serviceDevice.Service, logger log.Logger, resp *pool.Message, devices *coapSync.Map[uuid.UUID, *device]) error {
	dev, err := processDeviceResourceResponse(serviceDevice, logger, resp)
	if err != nil {
		return err
	}
	d, loaded := devices.LoadOrStore(dev.ID, dev)
	if !loaded {
		return nil
	}
	d.updateDeviceResourceBody(dev.private.DeviceResourceBody)

	return nil
}

func discoverDeviceResource(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onResponse func(conn *client.Conn, resp *pool.Message)) error {
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

func getDevicesByMulticast(ctx context.Context, discoveryCfg core.DiscoveryConfiguration, onDeviceResourceResponse, onDiscoveryResourceResponse func(conn *client.Conn, resp *pool.Message)) {
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

func getDeviceByMulticastAddress(ctx context.Context, serviceDevice *serviceDevice.Service, logger log.Logger, addr pkgNet.Addr, devices *coapSync.Map[uuid.UUID, *device]) error {
	hostname := addr.GetHostname()
	if strings.Contains(hostname, ":") {
		hostname = "[" + hostname + "]"
	}
	address := fmt.Sprintf("%s:%d", hostname, addr.GetPort())
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
	getDevicesByMulticast(ctx, discoveryConfiguration, func(conn *client.Conn, resp *pool.Message) {
		_ = conn.Close()
		err := onDeviceResourceResponse(serviceDevice, logger, resp, devices)
		if err != nil {
			return
		}
		if deviceResource.CompareAndSwap(false, true) && discoveryResource.Load() {
			cancel()
		}
	}, func(conn *client.Conn, resp *pool.Message) {
		_ = conn.Close()
		err := onDiscoveryResourceResponse(serviceDevice, logger, conn.RemoteAddr(), resp, devices)
		if err != nil {
			return
		}
		if discoveryResource.CompareAndSwap(false, true) && deviceResource.Load() {
			cancel()
		}
	})
	return nil
}

func getDeviceByAddress(ctx context.Context, serviceDevice *serviceDevice.Service, logger log.Logger, addr pkgNet.Addr, devices *coapSync.Map[uuid.UUID, *device]) error {
	if addr.GetPort() == MulticastPort {
		return getDeviceByMulticastAddress(ctx, serviceDevice, logger, addr, devices)
	}
	hostname := addr.GetHostname()
	if strings.Contains(hostname, ":") {
		hostname = "[" + hostname + "]"
	}
	address := fmt.Sprintf("%s:%d", hostname, addr.GetPort())
	client, err := udp.Dial(address, options.WithContext(ctx))
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
	deviceRes, err := processDeviceResourceResponse(serviceDevice, logger, resp)
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
	discoveryRes, err := processDiscoveryResourceResponse(serviceDevice, logger, client.RemoteAddr(), resp)
	if err != nil {
		return err
	}
	d, loaded := devices.LoadOrStore(deviceRes.ID, deviceRes)
	if loaded {
		d.updateDeviceResourceBody(deviceRes.private.DeviceResourceBody)
	}
	d.updateDeviceMetadata(discoveryRes.private.ResourceTypes, discoveryRes.private.Endpoints, discoveryRes.private.OwnershipStatus)
	return nil
}

func getDevicesByEndpoints(ctx context.Context, serviceDevice *serviceDevice.Service, logger log.Logger, endpoints []string, devices *coapSync.Map[uuid.UUID, *device]) {
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
		err := getDeviceByAddress(ctx, serviceDevice, logger, addresses[0], devices)
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
			err := getDeviceByAddress(ctx, serviceDevice, logger, addr, devices)
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

func (s *ClientApplicationServer) processDiscoverdDevices(discoveredDevices, cachedDevices *coapSync.Map[uuid.UUID, *device]) devices {
	devs := make(devices, 0, 128)
	discoveredDevices.Range(func(key uuid.UUID, d *device) bool {
		if len(d.GetEndpoints()) == 0 {
			// we don't want to return devices with no endpoints
			return true
		}

		updDevice, loaded := s.devices.LoadOrStore(key, d)
		if loaded {
			updDevice.update(d)
		}
		devs = append(devs, d)
		cachedDevices.Delete(key)
		return true
	})
	return devs
}

func sendDevices(req *pb.GetDevicesRequest, devs devices, send func(*grpcgwPb.Device) error) error {
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
		if err := send(d); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientApplicationServer) GetDevices(req *pb.GetDevicesRequest, srv pb.ClientApplication_GetDevicesServer) error {
	req = tryToSetDefaultRequest(req)
	ctx := srv.Context()
	var toCall []func()
	discoveredDevices := coapSync.NewMap[uuid.UUID, *device]()
	cachedDevices := coapSync.NewMap[uuid.UUID, *device]()
	timeout := DefaultTimeout
	if req.GetTimeout() > 0 {
		timeout = time.Duration(req.GetTimeout()) * time.Nanosecond
	}
	discoveryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if req.UseCache {
		s.devices.Range(func(key uuid.UUID, value *device) bool {
			cachedDevices.Store(key, value)
			return true
		})
	}
	if len(req.GetUseMulticast()) > 0 {
		toCall = append(toCall, func() {
			getDevicesByMulticast(discoveryCtx, toDiscoveryConfiguration(toUseMulticastFilter(req.GetUseMulticast())), func(conn *client.Conn, resp *pool.Message) {
				_ = conn.Close()
				_ = onDeviceResourceResponse(s.serviceDevice, s.logger, resp, discoveredDevices)
			}, func(conn *client.Conn, resp *pool.Message) {
				_ = conn.Close()
				_ = onDiscoveryResourceResponse(s.serviceDevice, s.logger, conn.RemoteAddr(), resp, discoveredDevices)
			})
		},
		)
	}
	if len(req.GetUseEndpoints()) > 0 {
		toCall = append(toCall, func() {
			getDevicesByEndpoints(discoveryCtx, s.serviceDevice, s.logger, req.GetUseEndpoints(), discoveredDevices)
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

	devs := s.processDiscoverdDevices(discoveredDevices, cachedDevices)
	cachedDevices.Range(func(key uuid.UUID, d *device) bool {
		devs = append(devs, d)
		return true
	})
	return sendDevices(req, devs, srv.Send)
}
