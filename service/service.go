//***************************************************************************
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
//*************************************************************************

package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	config         Config
	ctx            context.Context
	cancel         context.CancelFunc
	logger         log.Logger
	httpService    *http.Service
	grpcService    *grpc.Service
	deviceService  *device.Service
	tracerProvider trace.TracerProvider
	sigs           chan os.Signal
}

const serviceName = "client-application"

// New creates server.
func New(ctx context.Context, config Config, logger log.Logger) (*Service, error) {
	tracerProvider := trace.NewNoopTracerProvider()
	var httpService *http.Service
	deviceService, err := device.New(ctx, serviceName, config.Clients.Device, logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("cannot create device service: %w", err)
	}

	clientApplicationServer := grpc.NewClientApplicationServer(deviceService, logger)

	if config.APIs.HTTP.Enabled {
		httpService, err = http.New(ctx, serviceName, config.APIs.HTTP.Config, clientApplicationServer, logger, tracerProvider)
		if err != nil {
			return nil, fmt.Errorf("cannot create http service: %w", err)
		}
	}
	var grpcService *grpc.Service
	if config.APIs.GRPC.Enabled {
		grpcService, err = grpc.New(ctx, serviceName, config.APIs.GRPC.Config, clientApplicationServer, logger, tracerProvider)
		if err != nil {
			return nil, fmt.Errorf("cannot create grpc service: %w", err)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	s := Service{
		config: config,

		sigs: make(chan os.Signal, 1),

		ctx:    ctx,
		cancel: cancel,

		logger:         logger,
		httpService:    httpService,
		grpcService:    grpcService,
		deviceService:  deviceService,
		tracerProvider: tracerProvider,
	}

	return &s, nil
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

// Serve starts a device provisioning service on the configured address in *Service.
func (server *Service) Serve() error {
	if server.httpService != nil {
		scheme := "http"
		if server.config.APIs.HTTP.Config.TLS.Enabled {
			scheme = "https"
		}
		addr := getAddress(server.httpService.Address())
		log.Infof("HTTP API available on %v://%s%v", scheme, addr, http.ApiV1)
		if server.config.APIs.HTTP.UI.Enabled {
			log.Infof("HTTP UI available on %v://%s", scheme, addr)
		}
	}
	if server.grpcService != nil {
		addr := getAddress(server.grpcService.Address())
		insecure := ""
		if !server.config.APIs.GRPC.Config.TLS.Enabled {
			insecure = " (insecure)"
		}
		log.Infof("gRPC API available on %s%v", addr, insecure)
	}

	return server.serveWithHandlingSignal()
}

func (server *Service) serveWithHandlingSignal() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	services := make([]func() error, 0, 3)
	services = append(services, server.deviceService.Serve)
	if server.httpService != nil {
		services = append(services, server.httpService.Serve)
	}
	if server.grpcService != nil {
		services = append(services, server.grpcService.Serve)
	}
	wg.Add(len(services))
	for _, serve := range services {
		go func(serve func() error) {
			defer wg.Done()
			err := serve()
			errCh <- err
		}(serve)
	}

	signal.Notify(server.sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-server.sigs
	if server.httpService != nil {
		err := server.httpService.Shutdown()
		errCh <- err
	}
	if server.grpcService != nil {
		server.grpcService.Stop()
	}
	server.deviceService.Close()
	wg.Wait()

	var errors []error
	for {
		select {
		case err := <-errCh:
			if err != nil {
				errors = append(errors, err)
			}
		default:
			switch len(errors) {
			case 0:
				return nil
			case 1:
				return errors[0]
			default:
				return fmt.Errorf("%v", errors)
			}
		}
	}
}

// Shutdown turn off server.
func (server *Service) Close() error {
	select {
	case server.sigs <- syscall.SIGTERM:
	default:
	}
	return nil
}
