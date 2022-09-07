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
	"fmt"
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
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetHttp "github.com/plgd-dev/hub/v2/pkg/net/http"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	httpServer *http.Server
	listener   listener.Listener
}

type RequestHandler struct {
	mux                     *runtime.ServeMux
	clientApplicationServer *grpc.ClientApplicationServer
	config                  Config
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

// New creates new HTTP service
func New(ctx context.Context, serviceName string, config Config, clientApplicationServer *grpc.ClientApplicationServer, fileWatcher *fsnotify.Watcher, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	listener, err := listener.New(config.Config, fileWatcher, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	ch := new(inprocgrpc.Channel)
	pb.RegisterClientApplicationServer(ch, clientApplicationServer)
	grpcClient := pb.NewClientApplicationClient(ch)

	auth := func(ctx context.Context, method, uri string) (context.Context, error) {
		return ctx, nil
	}
	if clientApplicationServer.HasAuthorizationEnabled() {
		whiteList := []kitNetHttp.RequestMatcher{
			{
				Method: http.MethodGet,
				URI:    regexp.MustCompile(regexp.QuoteMeta(WebConfiguration)),
			},
			{
				Method: http.MethodGet,
				URI:    regexp.MustCompile(regexp.QuoteMeta(WellKnownJWKs)),
			},
			{
				// token is directly verified by clientApplication
				Method: http.MethodPut,
				URI:    regexp.MustCompile(regexp.QuoteMeta(WellKnownJWKs)),
			},
		}
		if config.UI.Enabled {
			whiteList = append(whiteList, kitNetHttp.RequestMatcher{
				Method: http.MethodGet,
				URI:    regexp.MustCompile(`^\/(a$|[^a].*|ap$|a[^p].*|ap[^i].*|api[^/])`),
			})
		}
		auth = kitNetHttp.NewInterceptorWithValidator(clientApplicationServer, authRules, whiteList...)
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
	if err := pb.RegisterClientApplicationHandlerClient(context.Background(), mux, grpcClient); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to register grpc-gateway handler: %w", err)
	}
	requestHandler := &RequestHandler{mux: mux, clientApplicationServer: clientApplicationServer, config: config}
	r.PathPrefix(Devices).Methods(http.MethodPut).MatcherFunc(resourceMatcher).HandlerFunc(requestHandler.updateResource)
	r.PathPrefix(Devices).Methods(http.MethodPost).MatcherFunc(resourceMatcher).HandlerFunc(requestHandler.createResource)
	r.PathPrefix(ApiV1).Handler(mux)
	r.HandleFunc(WebConfiguration, requestHandler.getWebConfiguration).Methods(http.MethodGet)
	r.HandleFunc(WellKnownJWKs, requestHandler.getJSONWebKeys).Methods(http.MethodGet)
	r.HandleFunc(WellKnownJWKs, requestHandler.updateJSONWebKeys).Methods(http.MethodPut)
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
				log.Errorf("failed to write response body: %w", err)
			}
		}))
	}

	httpServer := &http.Server{
		Handler:           kitNetHttp.OpenTelemetryNewHandler(handler, serviceName, tracerProvider),
		ReadTimeout:       config.Server.ReadTimeout,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
	}

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

var authRules = map[string][]kitNetHttp.AuthArgs{
	http.MethodGet: {
		{
			URI: regexp.MustCompile(regexp.QuoteMeta(ApiV1) + `\/.*`),
		},
	},
	http.MethodPost: {
		{
			URI: regexp.MustCompile(regexp.QuoteMeta(ApiV1) + `\/.*`),
		},
	},
	http.MethodDelete: {
		{
			URI: regexp.MustCompile(regexp.QuoteMeta(ApiV1) + `\/.*`),
		},
	},
	http.MethodPut: {
		{
			URI: regexp.MustCompile(regexp.QuoteMeta(ApiV1) + `\/.*`),
		},
	},
}
