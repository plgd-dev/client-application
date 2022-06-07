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
	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/client-application/service"
	serviceGrpc "github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/test/config"
	"github.com/plgd-dev/kit/v2/codec/cbor"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	CLIENT_APPLICATION_HTTP_HOST = "localhost:40050"
)

var (
	DevsimName string
)

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

func MakeConfig(t require.TestingT) service.Config {
	var cfg service.Config
	cfg.Log = log.MakeDefaultConfig()
	cfg.APIs.HTTP = MakeHttpConfig()

	return cfg
}

func SetUp(t *testing.T) (TearDown func()) {
	return New(t, MakeConfig(t))
}

// New creates test coap-gateway.
func New(t *testing.T, cfg service.Config) func() {
	ctx := context.Background()
	logger := log.NewLogger(cfg.Log)

	s, err := service.New(ctx, cfg, logger)
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

func MakeHttpConfig() http.Config {
	cfg := config.MakeListenerConfig(CLIENT_APPLICATION_HTTP_HOST)
	return http.Config{
		Addr: cfg.Addr,
		TLS: listener.TLSConfig{
			Enabled: true,
			Config:  cfg.TLS,
		},
	}
}

func NewHttpService(ctx context.Context, t *testing.T) (*http.Service, func()) {
	cfg := MakeConfig(t)
	cfg.APIs.HTTP.TLS.ClientCertificateRequired = false
	logger := log.NewLogger(cfg.Log)

	s, err := http.New(ctx, "client-application-http", cfg.APIs.HTTP, logger, trace.NewNoopTracerProvider())
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.Serve()
	}()

	cleanUp := func() {
		err = s.Shutdown()
		require.NoError(t, err)
		wg.Wait()
	}

	return s, cleanUp
}

type DeviceGatewayGetDevicesServer struct {
	grpc.ServerStream
	Devices []*pb.Device
	Ctx     context.Context
}

func NewDeviceGatewayGetDevicesServer(ctx context.Context) *DeviceGatewayGetDevicesServer {
	return &DeviceGatewayGetDevicesServer{
		Ctx: ctx,
	}
}

func (s *DeviceGatewayGetDevicesServer) Send(d *pb.Device) error {
	s.Devices = append(s.Devices, d)
	return nil
}

func (s *DeviceGatewayGetDevicesServer) Context() context.Context {
	return s.Ctx
}

func FindDeviceByName(name string, useMulticast []pb.GetDevicesRequest_UseMulticast) (*pb.Device, error) {
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		srv := NewDeviceGatewayGetDevicesServer(ctx)
		s := &serviceGrpc.DeviceGatewayServer{}
		err := s.GetDevices(&pb.GetDevicesRequest{
			UseMulticast: useMulticast,
		}, srv)
		if err != nil {
			return nil, err
		}
		for _, d := range srv.Devices {
			var dev device.Device
			if err := cbor.Decode(d.GetContent().GetData(), &dev); err != nil {
				continue
			}
			if dev.Name == name {
				return d, nil
			}
		}
	}
	return nil, fmt.Errorf("device %s not found", name)
}

func MustFindDeviceByName(name string, useMulticast []pb.GetDevicesRequest_UseMulticast) *pb.Device {
	d, err := FindDeviceByName(name, useMulticast)
	if err != nil {
		panic(err)
	}
	return d
}
