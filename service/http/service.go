package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/gorilla/handlers"
	router "github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetHttp "github.com/plgd-dev/hub/v2/pkg/net/http"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	httpServer *http.Server
	listener   listener.Listener
}

type RequestHandler struct {
	mux *runtime.ServeMux
}

func splitURIPath(requestURI, prefix string) []string {
	v := kitNetHttp.CanonicalHref(requestURI)
	p := strings.TrimPrefix(v, prefix) // remove core prefix
	if p == v {
		return nil
	}
	p = strings.TrimLeft(p, "/")
	return strings.Split(p, "/")
}

func resourceMatcher(r *http.Request, rm *router.RouteMatch) bool {
	paths := splitURIPath(r.RequestURI, Devices)
	if len(paths) > 2 &&
		paths[1] == ResourcesPathKey {
		if rm.Vars == nil {
			rm.Vars = make(map[string]string)
		}
		rm.Vars[DeviceIDKey] = paths[0]
		rm.Vars[ResourceHrefKey] = strings.Split("/"+strings.Join(paths[2:], "/"), "?")[0]
		return true
	}
	return false
}

// New creates new HTTP service
func New(ctx context.Context, serviceName string, config Config, deviceGatewayServer *grpc.DeviceGatewayServer, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	listener, err := listener.New(config.Config, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	ch := new(inprocgrpc.Channel)
	pb.RegisterDeviceGatewayServer(ch, deviceGatewayServer)
	grpcClient := pb.NewDeviceGatewayClient(ch)

	auth := func(ctx context.Context, method, uri string) (context.Context, error) {
		return ctx, nil
	}
	mux := serverMux.New()
	r := serverMux.NewRouter(queryCaseInsensitive, auth)

	corsOptions := make([]handlers.CORSOption, 0, 5)
	corsOptions = append(corsOptions, handlers.AllowedHeaders(config.CORS.AllowedHeaders))
	corsOptions = append(corsOptions, handlers.AllowedOrigins(config.CORS.AllowedOrigins))
	corsOptions = append(corsOptions, handlers.AllowedMethods(config.CORS.AllowedMethods))
	if config.CORS.AllowCredentials {
		corsOptions = append(corsOptions, handlers.AllowCredentials())
	}
	handler := handlers.CORS(corsOptions...)(r)

	// register grpc-proxy handler
	if err := pb.RegisterDeviceGatewayHandlerClient(context.Background(), mux, grpcClient); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to register grpc-gateway handler: %w", err)
	}
	requestHandler := &RequestHandler{mux: mux}
	r.PathPrefix(Devices).Methods(http.MethodPut).MatcherFunc(resourceMatcher).HandlerFunc(requestHandler.updateResource)
	r.PathPrefix("/").Handler(mux)
	httpServer := &http.Server{Handler: kitNetHttp.OpenTelemetryNewHandler(handler, serviceName, tracerProvider)}

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
	UseCacheQueryKey              = "useCache"
	UseMulticastQueryKey          = "useMulticast"
	UseEndpointsQueryKey          = "useEndpoints"
	TimeoutQueryKey               = "timeout"
	OwnershipStatusFilterQueryKey = "ownershipStatusFilter"
	TypeFilterQueryKey            = "typeFilter"
)

var queryCaseInsensitive = map[string]string{
	strings.ToLower(UseCacheQueryKey):              UseCacheQueryKey,
	strings.ToLower(UseMulticastQueryKey):          UseMulticastQueryKey,
	strings.ToLower(UseEndpointsQueryKey):          UseEndpointsQueryKey,
	strings.ToLower(TimeoutQueryKey):               TimeoutQueryKey,
	strings.ToLower(OwnershipStatusFilterQueryKey): OwnershipStatusFilterQueryKey,
	strings.ToLower(TypeFilterQueryKey):            TypeFilterQueryKey,
}
