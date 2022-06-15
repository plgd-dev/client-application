package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	pkgGrpcServer "github.com/plgd-dev/client-application/pkg/net/grpc/server"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	grpcServer *server.Server
}

// New creates new GRPC service
func New(ctx context.Context, serviceName string, config Config, deviceGatewayServer *DeviceGatewayServer, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	interceptor := kitNetGrpc.MakeAuthInterceptors(func(ctx context.Context, method string) (context.Context, error) {
		return ctx, nil
	})
	opts, err := server.MakeDefaultOptions(interceptor, logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server options: %w", err)
	}

	server, err := pkgGrpcServer.New(config, logger, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}
	pb.RegisterDeviceGatewayServer(server.Server, deviceGatewayServer)

	return &Service{
		grpcServer: server,
	}, nil
}

// Serve starts the service's HTTP server and blocks
func (s *Service) Serve() error {
	return s.grpcServer.Serve()
}

// Shutdown ends serving
func (s *Service) Stop() {
	s.grpcServer.Stop()
}

func (s *Service) Address() string {
	return s.grpcServer.Addr()
}
