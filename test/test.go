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

package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	_ "cloud.google.com/go"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/net/grpc/server"
	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/client-application/service/config"
	configDevice "github.com/plgd-dev/client-application/service/config/device"
	configGrpc "github.com/plgd-dev/client-application/service/config/grpc"
	configHttp "github.com/plgd-dev/client-application/service/config/http"
	"github.com/plgd-dev/client-application/service/config/remoteProvisioning"
	serviceDevice "github.com/plgd-dev/client-application/service/device"
	serviceGrpc "github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/device/v2/schema"
	"github.com/plgd-dev/device/v2/schema/device"
	deviceTest "github.com/plgd-dev/device/v2/test"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	testConfig "github.com/plgd-dev/hub/v2/test/config"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	CLIENT_APPLICATION_HTTP_HOST = "localhost:40050"
	CLIENT_APPLICATION_GRPC_HOST = "localhost:40051"
	VERSION                      = "v0.0.1-test"
	BUILD_DATE                   = "1.1.1970"
	COMMIT_HASH                  = "aaa"
	PSK_OWNER                    = "57b3fae9-adf5-4e34-90ea-e77784407103"
)

var (
	MFG_ROOT_CA_CRT            = os.Getenv("MFG_ROOT_CA_CRT")
	MFG_CLIENT_APPLICATION_CRT = os.Getenv("MFG_CLIENT_APPLICATION_CRT")
	MFG_CLIENT_APPLICATION_KEY = os.Getenv("MFG_CLIENT_APPLICATION_KEY")
)

var DevsimName string

func init() {
	DevsimName = "devsim-" + MustGetHostname()
}

func MustGetHostname() string {
	n, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return n
}

func MakeConfig(t require.TestingT) config.Config {
	var cfg config.Config
	cfg.Log = log.MakeDefaultConfig()
	cfg.APIs.HTTP = MakeHttpConfig()
	cfg.APIs.GRPC = MakeGrpcConfig()
	cfg.Clients.Device = MakeDeviceConfig()
	cfg.RemoteProvisioning = MakeRemoteProvisioningConfig()

	require.NoError(t, cfg.Validate())

	return cfg
}

func SetUp(t *testing.T) (tearDown func()) {
	return New(t, MakeConfig(t))
}

// New creates test coap-gateway.
func New(t *testing.T, cfg config.Config) func() {
	ctx := context.Background()
	logger := log.NewLogger(cfg.Log)
	require.NoError(t, cfg.Validate())

	fileWatcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	s, err := service.New(ctx, cfg, NewServiceInformation(), fileWatcher, logger)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.Serve()
	}()

	return func() {
		_ = s.Close()
		wg.Wait()
	}
}

func MakeDeviceConfig() configDevice.Config {
	cfg := configDevice.Config{
		COAP: configDevice.CoapConfig{
			MaxMessageSize: 256 * 1024,
			InactivityMonitor: configDevice.InactivityMonitor{
				Timeout: time.Second * 1,
			},
			BlockwiseTransfer: configDevice.BlockwiseTransferConfig{
				Enabled: true,
				SZXStr:  "1024",
			},
			TLS: configDevice.TLSConfig{
				Authentication: configDevice.AuthenticationPreSharedKey,
				PreSharedKey: configDevice.PreSharedKeyConfig{
					SubjectUUIDStr: PSK_OWNER,
					KeyUUIDStr:     "46178d21-d480-4e95-9bd3-6c9eefa8d9d8",
				},
			},
			OwnershipTransfer: configDevice.OwnershipTransferConfig{
				Methods: []configDevice.OwnershipTransferMethod{configDevice.OwnershipTransferJustWorks, configDevice.OwnershipTransferManufacturerCertificate},
				Manufacturer: configDevice.ManufacturerConfig{
					TLS: configDevice.ManufacturerTLSConfig{
						CAPool:   MFG_ROOT_CA_CRT,
						CertFile: MFG_CLIENT_APPLICATION_CRT,
						KeyFile:  MFG_CLIENT_APPLICATION_KEY,
					},
				},
			},
		},
	}
	return cfg
}

func MakeHttpConfig() config.HTTPConfig {
	cfg := testConfig.MakeListenerConfig(CLIENT_APPLICATION_HTTP_HOST)
	return config.HTTPConfig{
		Enabled: true,
		Config: configHttp.Config{
			Config: listener.Config{
				Addr: cfg.Addr,
				TLS: listener.TLSConfig{
					Enabled: true,
					Config:  cfg.TLS,
				},
			},
			CORS: configHttp.CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedHeaders: []string{"Accept", "Accept-Language", "Accept-Encoding", "Content-Type", "Content-Language", "Content-Length", "Origin", "X-CSRF-Token", "Authorization"},
				AllowedMethods: []string{"GET", "PATCH", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
			},
		},
	}
}

func MakeGrpcConfig() config.GRPCConfig {
	cfg := testConfig.MakeGrpcServerConfig(CLIENT_APPLICATION_GRPC_HOST)
	return config.GRPCConfig{
		Enabled: true,
		Config: configGrpc.Config{
			Addr:              cfg.Addr,
			EnforcementPolicy: cfg.EnforcementPolicy,
			KeepAlive:         cfg.KeepAlive,
			TLS: server.TLSConfig{
				Enabled: true,
				Config:  cfg.TLS,
			},
		},
	}
}

func MakeRemoteProvisioningConfig() remoteProvisioning.Config {
	return remoteProvisioning.Config{
		Mode: remoteProvisioning.Mode_None,
		UserAgentConfig: remoteProvisioning.UserAgentConfig{
			CSRChallengeStateExpiration: time.Minute * 10,
			CertificateAuthorityAddress: testConfig.CERTIFICATE_AUTHORITY_HOST,
		},
		Authorization: remoteProvisioning.AuthorizationConfig{
			OwnerClaim: testConfig.OWNER_CLAIM,
			Authority:  testConfig.OAUTH_SERVER_HOST,
			ClientID:   testConfig.OAUTH_MANAGER_CLIENT_ID,
		},
	}
}

func NewHttpService(ctx context.Context, t *testing.T) (*http.Service, func()) {
	cfg := MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	logger := log.NewLogger(cfg.Log)
	clientApplicationServer, tearDown, err := NewClientApplicationServer(ctx)
	require.NoError(t, err)

	fileWatcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	s, err := http.New(ctx, "client-application-http", cfg.APIs.HTTP.Config, clientApplicationServer, fileWatcher, logger, trace.NewNoopTracerProvider())
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.Serve()
	}()

	cleanUp := func() {
		err = s.Close()
		require.NoError(t, err)
		wg.Wait()
		tearDown()
	}

	return s, cleanUp
}

func NewServiceInformation() *configGrpc.ServiceInformation {
	return &configGrpc.ServiceInformation{
		Version:                  VERSION,
		BuildDate:                BUILD_DATE,
		CommitHash:               COMMIT_HASH,
		IsInitialized:            true,
		DeviceAuthenticationMode: pb.GetConfigurationResponse_PRE_SHARED_KEY,
		Owner:                    PSK_OWNER,
	}
}

type ClientApplicationServerCfg struct {
	Cfg                   configDevice.Config
	RemoteProvisioningCfg remoteProvisioning.Config
}

type ClientApplicationServerOpt = func(c *ClientApplicationServerCfg)

func WithDeviceConfig(cfg configDevice.Config) ClientApplicationServerOpt {
	return func(c *ClientApplicationServerCfg) {
		c.Cfg = cfg
	}
}

func WithRemoteProvisioningConfig(cfg remoteProvisioning.Config) ClientApplicationServerOpt {
	return func(c *ClientApplicationServerCfg) {
		c.RemoteProvisioningCfg = cfg
	}
}

func NewClientApplicationServer(ctx context.Context, opts ...ClientApplicationServerOpt) (*serviceGrpc.ClientApplicationServer, func(), error) {
	logger := log.NewLogger(log.MakeDefaultConfig())
	cfg := ClientApplicationServerCfg{
		Cfg:                   MakeDeviceConfig(),
		RemoteProvisioningCfg: MakeRemoteProvisioningConfig(),
	}
	for _, o := range opts {
		o(&cfg)
	}
	deviceCfg := cfg.Cfg
	if err := deviceCfg.Validate(); err != nil {
		return nil, nil, err
	}
	remoteProvisioningCfg := cfg.RemoteProvisioningCfg
	if err := remoteProvisioningCfg.Validate(); err != nil {
		return nil, nil, err
	}
	d, err := serviceDevice.New(ctx, "client-application-device", deviceCfg, logger, trace.NewNoopTracerProvider())
	if err != nil {
		return nil, nil, err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = d.Serve()
	}()

	clientApplicationServer := serviceGrpc.NewClientApplicationServer(remoteProvisioningCfg, d, NewServiceInformation(), logger)
	return clientApplicationServer, func() {
		_ = d.Close()
		clientApplicationServer.Close()
		wg.Wait()
	}, nil
}

type ClientApplicationGetDevicesServer struct {
	grpc.ServerStream
	Devices []*grpcgwPb.Device
	Ctx     context.Context
}

func NewClientApplicationGetDevicesServer(ctx context.Context) *ClientApplicationGetDevicesServer {
	return &ClientApplicationGetDevicesServer{
		Ctx: ctx,
	}
}

func (s *ClientApplicationGetDevicesServer) Send(d *grpcgwPb.Device) error {
	s.Devices = append(s.Devices, d)
	return nil
}

func (s *ClientApplicationGetDevicesServer) Context() context.Context {
	return s.Ctx
}

func FindDeviceByName(name string, useMulticast []pb.GetDevicesRequest_UseMulticast) (*grpcgwPb.Device, error) {
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		srv := NewClientApplicationGetDevicesServer(ctx)
		s, teardown, err := NewClientApplicationServer(ctx)
		if err != nil {
			return nil, err
		}
		defer teardown()
		err = s.GetDevices(&pb.GetDevicesRequest{
			UseMulticast: useMulticast,
		}, srv)
		if err != nil {
			return nil, err
		}
		for _, d := range srv.Devices {
			var dev device.Device
			if err := cbor.Decode(d.GetData().GetContent().GetData(), &dev); err != nil {
				continue
			}
			if dev.Name == name {
				return d, nil
			}
		}
	}
	return nil, fmt.Errorf("device %s not found", name)
}

func MustFindDeviceByName(name string, useMulticast []pb.GetDevicesRequest_UseMulticast) *grpcgwPb.Device {
	d, err := FindDeviceByName(name, useMulticast)
	if err != nil {
		panic(err)
	}
	return d
}

func GetDeviceResourceLinks() schema.ResourceLinks {
	resources := make(schema.ResourceLinks, 0, len(deviceTest.TestDevsimResources)+len(deviceTest.TestDevsimPrivateResources)+len(deviceTest.TestDevsimSecResources))
	resources = append(resources, deviceTest.TestDevsimResources...)
	resources = append(resources, deviceTest.TestDevsimPrivateResources...)
	resources = append(resources, deviceTest.TestDevsimSecResources...)
	return resources
}
