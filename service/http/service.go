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

package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/gorilla/handlers"
	router "github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/client-application/pkg/net/listener/tls"
	configHttp "github.com/plgd-dev/client-application/service/config/http"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	pkgLog "github.com/plgd-dev/hub/v2/pkg/log"
	kitNetHttp "github.com/plgd-dev/hub/v2/pkg/net/http"
	pkgHttpJwt "github.com/plgd-dev/hub/v2/pkg/net/http/jwt"
	pkgHttpUri "github.com/plgd-dev/hub/v2/pkg/net/http/uri"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	httpServer *http.Server
	listener   listener.Listener
}

type RequestHandler struct {
	mux                     *runtime.ServeMux
	clientApplicationServer *grpc.ClientApplicationServer
	config                  configHttp.Config
}

func splitURIPath(requestURI, prefix string) []string {
	v := pkgHttpUri.CanonicalHref(requestURI)
	p := strings.TrimPrefix(v, prefix) // remove core prefix
	if p == v {
		return nil
	}
	p = strings.TrimLeft(p, "/")
	return strings.Split(p, "/")
}

func resourceMatcher(r *http.Request, rm *router.RouteMatch) bool {
	paths := splitURIPath(r.RequestURI, Devices)
	if len(paths) > 2 && (paths[1] == ResourcesPathKey || paths[1] == ResourceLinksPathKey) {
		if rm.Vars == nil {
			rm.Vars = make(map[string]string)
		}
		rm.Vars[DeviceIDKey] = paths[0]
		rm.Vars[ResourceHrefKey] = strings.Split("/"+strings.Join(paths[2:], "/"), "?")[0]
		return true
	}
	return false
}

func createAuthFunc(config configHttp.Config, clientApplicationServer *grpc.ClientApplicationServer) func(ctx context.Context, method, uri string) (context.Context, error) {
	whiteList := []pkgHttpJwt.RequestMatcher{
		{
			Method: http.MethodGet,
			URI:    regexp.MustCompile(regexp.QuoteMeta(WellKnownJWKs)),
		},
		{
			Method: http.MethodGet,
			URI:    regexp.MustCompile(regexp.QuoteMeta(WellKnownConfiguration)),
		},
		{
			// token is directly verified by clientApplication
			Method: http.MethodPost,
			URI:    regexp.MustCompile(regexp.QuoteMeta(Initialize)),
		},
	}
	if config.UI.Enabled {
		whiteList = append(whiteList, pkgHttpJwt.RequestMatcher{
			Method: http.MethodGet,
			URI:    regexp.MustCompile(`^\/(a$|[^a].*|ap$|a[^p].*|ap[^i].*|api[^/])`),
		})
	}
	auth := pkgHttpJwt.NewInterceptorWithValidator(clientApplicationServer, kitNetHttp.NewDefaultAuthorizationRules(ApiV1), whiteList...)
	return func(ctx context.Context, method, uri string) (context.Context, error) {
		if clientApplicationServer.HasJWTAuthorizationEnabled() {
			return auth(ctx, method, uri)
		}
		return ctx, nil
	}
}

type contextKey struct {
	key string
}

var connContextKey = &contextKey{"http-conn"}

func saveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, connContextKey, c)
}

func getTLSConn(r *http.Request) (*tls.Conn, bool) {
	c, ok := r.Context().Value(connContextKey).(*tls.Conn)
	return c, ok
}

func newListener(config configHttp.Config, fileWatcher *fsnotify.Watcher, logger pkgLog.Logger) (listener.Listener, error) {
	if config.Config.TLS.Enabled {
		return tls.New(tls.Config{
			Addr: config.Config.Addr,
			TLS:  config.Config.TLS.Config,
		}, fileWatcher, logger)
	}
	return listener.New(config.Config, fileWatcher, logger)
}

func wrapHandler(handler http.Handler, serviceName string, tracerProvider trace.TracerProvider) http.Handler {
	h := kitNetHttp.OpenTelemetryNewHandler(handler, serviceName, tracerProvider)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := getTLSConn(r); ok && c.ConnectionType == tls.ConnectionTypeHTTP {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			_ = c.Close()
			return
		}
		h.ServeHTTP(w, r)
	})
}

func newCORSHandler(config configHttp.Config, r *router.Router) http.Handler {
	corsOptions := make([]handlers.CORSOption, 0, 5)
	corsOptions = append(corsOptions, handlers.AllowedHeaders(config.CORS.AllowedHeaders))
	corsOptions = append(corsOptions, handlers.AllowedOrigins(config.CORS.AllowedOrigins))
	corsOptions = append(corsOptions, handlers.AllowedMethods(config.CORS.AllowedMethods))
	if config.CORS.AllowCredentials {
		corsOptions = append(corsOptions, handlers.AllowCredentials())
	}
	return handlers.CORS(corsOptions...)(r)
}

func setUIHandlers(config configHttp.Config, r *router.Router) {
	// serve www directory
	if config.UI.Enabled {
		fs := http.FileServer(http.Dir(config.UI.Directory))
		r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := httptest.NewRecorder()
			fs.ServeHTTP(c, r)
			if c.Code == http.StatusNotFound {
				c = httptest.NewRecorder()
				r.URL.Path = "/"
				fs.ServeHTTP(c, r)
			}
			for k, v := range c.Header() {
				w.Header().Set(k, strings.Join(v, ""))
			}
			w.WriteHeader(c.Code)
			if _, err := c.Body.WriteTo(w); err != nil {
				pkgLog.Errorf("failed to write response body: %w", err)
			}
		}))
	}
}

type logWriter struct {
	logger pkgLog.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.logger.Debug(string(p))
	return len(p), nil
}

func newErrorLogger(logger pkgLog.Logger) *log.Logger {
	return log.New(&logWriter{logger: logger.With("package", "net/http")}, "", 0)
}

// New creates new HTTP service
func New(ctx context.Context, serviceName string, config configHttp.Config, clientApplicationServer *grpc.ClientApplicationServer, fileWatcher *fsnotify.Watcher, logger pkgLog.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	lis, err := newListener(config, fileWatcher, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	ch := new(inprocgrpc.Channel)
	pb.RegisterClientApplicationServer(ch, clientApplicationServer)
	grpcClient := pb.NewClientApplicationClient(ch)

	auth := createAuthFunc(config, clientApplicationServer)
	mux := serverMux.New()
	r := serverMux.NewRouter(queryCaseInsensitive, auth)
	handler := newCORSHandler(config, r)

	// register grpc-proxy handler
	if err := pb.RegisterClientApplicationHandlerClient(ctx, mux, grpcClient); err != nil {
		_ = lis.Close()
		return nil, fmt.Errorf("failed to register grpc-gateway handler: %w", err)
	}
	requestHandler := &RequestHandler{mux: mux, clientApplicationServer: clientApplicationServer, config: config}
	r.PathPrefix(Devices).Methods(http.MethodPut).MatcherFunc(resourceMatcher).HandlerFunc(requestHandler.updateResource)
	r.PathPrefix(Devices).Methods(http.MethodPost).MatcherFunc(resourceMatcher).HandlerFunc(requestHandler.createResource)
	r.PathPrefix(ApiV1).Handler(mux)
	r.PathPrefix(WellKnown).Handler(mux)

	setUIHandlers(config, r)

	httpServer := &http.Server{
		Handler:           wrapHandler(handler, serviceName, tracerProvider),
		ReadTimeout:       config.Server.ReadTimeout,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
		ConnContext:       saveConnInContext,
		ErrorLog:          newErrorLogger(logger),
	}

	return &Service{
		httpServer: httpServer,
		listener:   lis,
	}, nil
}

// Serve starts the service's HTTP server and blocks
func (s *Service) Serve() error {
	err := s.httpServer.Serve(s.listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Close serving
func (s *Service) Close() error {
	return s.httpServer.Shutdown(context.Background())
}

func (s *Service) Address() string {
	return s.listener.Addr().String()
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
