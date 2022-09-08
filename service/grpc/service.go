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

	"github.com/plgd-dev/client-application/pb"
	pkgGrpcServer "github.com/plgd-dev/client-application/pkg/net/grpc/server"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	grpcServer *server.Server
}

// New creates new GRPC service
func New(ctx context.Context, serviceName string, config Config, clientApplicationServer *ClientApplicationServer, fileWatcher *fsnotify.Watcher, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	interceptor := kitNetGrpc.MakeAuthInterceptors(func(ctx context.Context, method string) (context.Context, error) {
		return ctx, nil
	})
	if clientApplicationServer.HasJWTAuthorizationEnabled() {
		methods := []string{
			"/" + pb.ClientApplication_ServiceDesc.ServiceName + "/UpdateJSONWebKeys",
			"/" + pb.ClientApplication_ServiceDesc.ServiceName + "/GetJSONWebKeys",
			"/" + pb.ClientApplication_ServiceDesc.ServiceName + "/GetConfiguration",
		}
		interceptor = server.NewAuth(clientApplicationServer, server.WithWhiteListedMethods(methods...))
	}
	opts, err := server.MakeDefaultOptions(interceptor, logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server options: %w", err)
	}

	server, err := pkgGrpcServer.New(config, fileWatcher, logger, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}
	pb.RegisterClientApplicationServer(server.Server, clientApplicationServer)

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
