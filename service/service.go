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

package service

import (
	"context"
	"fmt"
	"net"

	"github.com/plgd-dev/client-application/service/config"
	configDevice "github.com/plgd-dev/client-application/service/config/device"
	configGrpc "github.com/plgd-dev/client-application/service/config/grpc"
	"github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/fn"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/pkg/service"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/atomic"
)

const serviceName = "client-application"

func newHttpService(ctx context.Context, config config.Config, clientApplicationServer *grpc.ClientApplicationServer, fileWatcher *fsnotify.Watcher, logger log.Logger, tracerProvider trace.TracerProvider) (*http.Service, error) {
	httpService, err := http.New(ctx, serviceName, config.APIs.HTTP.Config, clientApplicationServer, fileWatcher, logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("cannot create http service: %w", err)
	}
	scheme := "https"
	addr := getAddress(httpService.Address())
	log.Infof("HTTP API available on %v://%s%v", scheme, addr, http.ApiV1)
	if config.APIs.HTTP.UI.Enabled {
		log.Infof("HTTP UI available on %v://%s", scheme, addr)
	}
	return httpService, err
}

func newGrpcService(ctx context.Context, config config.Config, clientApplicationServer *grpc.ClientApplicationServer, fileWatcher *fsnotify.Watcher, logger log.Logger, tracerProvider trace.TracerProvider) (*grpc.Service, error) {
	grpcService, err := grpc.New(ctx, serviceName, config.APIs.GRPC.Config, clientApplicationServer, fileWatcher, logger, tracerProvider)
	if err != nil {
		return nil, err
	}
	addr := getAddress(grpcService.Address())
	insecure := ""
	if !config.APIs.GRPC.Config.TLS.Enabled {
		insecure = " (insecure)"
	}
	log.Infof("gRPC API available on %s%v", addr, insecure)
	return grpcService, nil
}

func closeServicesOnError(err error, services []service.APIService) error {
	errors := []error{err}
	for _, s := range services {
		switch s.(type) {
		case *http.Service:
			err := s.Close()
			if err != nil {
				errors = append(errors, fmt.Errorf("cannot close http service: %w", err))
			}
		case *device.Service:
			err := s.Close()
			if err != nil {
				errors = append(errors, fmt.Errorf("cannot close device service: %w", err))
			}
		}
	}
	if len(errors) == 1 {
		return errors[0]
	}
	return fmt.Errorf("%v", errors)
}

// New creates server.
func New(ctx context.Context, cfg config.Config, info *configGrpc.ServiceInformation, fileWatcher *fsnotify.Watcher, logger log.Logger) (*service.Service, error) {
	tracerProvider := noop.NewTracerProvider()
	var closerFunc fn.FuncList
	config := atomic.NewPointer(&cfg)
	var deviceService *device.Service
	var err error
	if cfg.Clients.Device.COAP.TLS.Authentication != configDevice.AuthenticationUninitialized {
		deviceService, err = device.New(ctx, func() configDevice.Config {
			return config.Load().Clients.Device
		}, logger)
		if err != nil {
			return nil, fmt.Errorf("cannot create device service: %w", err)
		}
	}
	clientApplicationServer := grpc.NewClientApplicationServer(config, deviceService, info, logger)
	closerFunc.AddFunc(clientApplicationServer.Close)
	services := make([]service.APIService, 0, 2)
	if cfg.APIs.HTTP.Enabled {
		httpService, err := newHttpService(ctx, cfg, clientApplicationServer, fileWatcher, logger, tracerProvider)
		if err != nil {
			closerFunc.Execute()
			return nil, fmt.Errorf("cannot create http service: %w", err)
		}
		services = append(services, httpService)
	}
	if cfg.APIs.GRPC.Enabled {
		grpcService, err := newGrpcService(ctx, cfg, clientApplicationServer, fileWatcher, logger, tracerProvider)
		if err != nil {
			closerFunc.Execute()
			return nil, closeServicesOnError(fmt.Errorf("cannot create grpc service: %w", err), services)
		}
		services = append(services, grpcService)
	}
	s := service.New(services...)
	s.AddCloseFunc(closerFunc.Execute)
	return s, nil
}

func getAddress(addr string) string {
	hostname, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	ip := net.ParseIP(hostname)
	if ip == nil {
		return addr
	}
	if ip.IsUnspecified() {
		if ip.To4() == nil {
			return "[::1]:" + port
		}
		return "127.0.0.1:" + port
	}
	return addr
}
