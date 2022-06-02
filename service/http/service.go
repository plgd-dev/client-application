package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetHttp "github.com/plgd-dev/hub/v2/pkg/net/http"
	"github.com/plgd-dev/hub/v2/pkg/net/listener"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	httpServer *http.Server
	listener   *listener.Server
}

// New creates new HTTP service
func New(ctx context.Context, serviceName string, config Config, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	listener, err := listener.New(config.Connection, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	ch := new(inprocgrpc.Channel)
	pb.RegisterDeviceGatewayServer(ch, grpc.NewDeviceGatewayServer())
	grpcClient := pb.NewDeviceGatewayClient(ch)

	auth := func(ctx context.Context, method, uri string) (context.Context, error) {
		return ctx, nil
	}
	mux := serverMux.New()
	r := serverMux.NewRouter(queryCaseInsensitive, auth)

	// register grpc-proxy handler
	if err := pb.RegisterDeviceGatewayHandlerClient(context.Background(), mux, grpcClient); err != nil {
		return nil, fmt.Errorf("failed to register grpc-gateway handler: %w", err)
	}
	r.PathPrefix("/").Handler(mux)
	httpServer := &http.Server{Handler: kitNetHttp.OpenTelemetryNewHandler(r, serviceName, tracerProvider)}

	return &Service{
		httpServer: httpServer,
		listener:   listener,
	}, nil
}

// Serve starts the service's HTTP server and blocks
func (s *Service) Serve() error {
	return s.httpServer.Serve(s.listener)
}

// Shutdown ends serving
func (s *Service) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}

const (
	IDFilterQueryKey                = "idFilter"
	EnrollmentGroupIdFilterQueryKey = "enrollmentGroupIdFilter"
	DeviceIdFilterQueryKey          = "deviceIdFilter"
)

var queryCaseInsensitive = map[string]string{
	strings.ToLower(IDFilterQueryKey):                IDFilterQueryKey,
	strings.ToLower(EnrollmentGroupIdFilterQueryKey): EnrollmentGroupIdFilterQueryKey,
	strings.ToLower(DeviceIdFilterQueryKey):          DeviceIdFilterQueryKey,
}
