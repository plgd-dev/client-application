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

package device

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/client/core/otm"
	justworks "github.com/plgd-dev/device/client/core/otm/just-works"
	"github.com/plgd-dev/device/client/core/otm/manufacturer"
	"github.com/plgd-dev/device/pkg/net/coap"
	coapNet "github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/net/monitor/inactivity"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type AuthenticationClient interface {
	DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error)
	DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error)
	GetOwnerID() (string, error)
	GetOwnOptions() []core.OwnOption

	GetIdentityCSR(id string) ([]byte, error)
	SetIdentityCertificate(chainPem []byte) error
	GetIdentityCertificate() (tls.Certificate, error)
	GetCertificateAuthorities() ([]*x509.Certificate, error)
	IsInitialized() bool
	Reset()
}

type Service struct {
	config               Config
	logger               log.Logger
	tracerProvider       trace.TracerProvider
	udp4server           *udp.Server
	udp6server           *udp.Server
	udp4Listener         *coapNet.UDPConn
	udp6Listener         *coapNet.UDPConn
	done                 chan struct{}
	authenticationClient AuthenticationClient
}

var closeHandlerKey = "close-handler"

// New creates new GRPC service
func New(ctx context.Context, serviceName string, config Config, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	var authenticationClient AuthenticationClient
	switch config.COAP.TLS.Authentication {
	case AuthenticationPreSharedKey:
		authenticationClient = newAuthenticationPreSharedKey(config)
	case AuthenticationX509:
		authenticationClient = newAuthenticationX509(config)
	}

	opts := []udp.ServerOption{
		udp.WithContext(ctx),
		udp.WithOnNewClientConn(func(cc *client.ClientConn) {
			closeHander := coap.NewOnCloseHandler()
			cc.AddOnClose(func() {
				closeHander.OnClose(nil)
			})
			cc.SetContextValue(&closeHandlerKey, closeHander)
		}),
		udp.WithMaxMessageSize(config.COAP.MaxMessageSize),
		udp.WithBlockwise(config.COAP.BlockwiseTransfer.Enabled, config.COAP.BlockwiseTransfer.szx, config.COAP.InactivityMonitor.Timeout),
		udp.WithErrors(errFunc),
		udp.WithInactivityMonitor(config.COAP.InactivityMonitor.Timeout, func(cc inactivity.ClientConn) {
			if c, ok := cc.(*client.ClientConn); ok {
				log.Debugf("closing connection %v for inactivity", c.RemoteAddr())
			}
			_ = cc.Close()
		}),
	}

	udp4Listener, err := coapNet.NewListenUDP("udp4", ":")
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP4 listener: %w", err)
	}

	udp6Listener, err := coapNet.NewListenUDP("udp6", ":")
	if err != nil {
		_ = udp4Listener.Close()
		return nil, fmt.Errorf("failed to create UDP6 listener: %w", err)
	}

	udp4server := udp.NewServer(opts...)
	udp6server := udp.NewServer(opts...)

	return &Service{
		config:               config,
		logger:               logger,
		tracerProvider:       tracerProvider,
		udp4server:           udp4server,
		udp6server:           udp6server,
		udp4Listener:         udp4Listener,
		udp6Listener:         udp6Listener,
		done:                 make(chan struct{}),
		authenticationClient: authenticationClient,
	}, nil
}

func errFunc(err error) {
	log.Debug(err)
}

func (s *Service) DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return s.authenticationClient.DialDTLS(ctx, addr, dtlsCfg, append(dialOpts, opts...)...)
}

func (s *Service) DialOwnership(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, append(dialOpts, opts...)...)
}

type UDPClientConn struct {
	*client.Client
}

func (c *UDPClientConn) Context() context.Context {
	if cc, ok := c.Client.ClientConn().(*client.ClientConn); ok {
		// we need to check if connection will be closed for inactivity
		cc.CheckExpirations(time.Now().Add(50 * time.Millisecond))
		if cc.Context().Err() == nil {
			// move inactivity timeout to future
			cc.InactivityMonitor().Notify()
		}
	}
	return c.Client.Context()
}

func (s *Service) DialUDP(ctx context.Context, addr string, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	var cc *client.ClientConn
	if udpAddr.IP.To4() != nil {
		cc, err = s.udp4server.NewClientConn(udpAddr)
	} else {
		cc, err = s.udp6server.NewClientConn(udpAddr)
	}
	if err != nil {
		return nil, err
	}
	closeHandler := cc.Context().Value(&closeHandlerKey)
	if closeHandler == nil {
		_ = cc.Close()
		return nil, fmt.Errorf("failed to create client connection: close handler is nil")
	}
	h, ok := closeHandler.(*coap.OnCloseHandler)
	if !ok {
		_ = cc.Close()
		return nil, fmt.Errorf("failed to create client connection: close handler is not *coap.OnCloseHandler")
	}
	return coap.NewClientCloseHandler(&UDPClientConn{Client: cc.Client()}, h), nil
}

func (s *Service) DialTCP(ctx context.Context, addr string, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return coap.DialTCP(ctx, addr, append(dialOpts, opts...)...)
}

func (s *Service) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return s.authenticationClient.DialTLS(ctx, addr, tlsCfg, append(dialOpts, opts...)...)
}

func (s *Service) ErrFunc(err error) {
	log.Debug(err)
}

func (s *Service) GetDeviceConfiguration() core.DeviceConfiguration {
	return core.DeviceConfiguration{
		DialDTLS: s.DialDTLS,
		DialUDP:  s.DialUDP,
		DialTCP:  s.DialTCP,
		DialTLS:  s.DialTLS,
		ErrFunc:  s.ErrFunc,
		TLSConfig: &core.TLSConfig{
			GetCertificate:            s.authenticationClient.GetIdentityCertificate,
			GetCertificateAuthorities: s.authenticationClient.GetCertificateAuthorities,
		},
		GetOwnerID: s.authenticationClient.GetOwnerID,
	}
}

func (s *Service) GetOwnershipClients() []otm.Client {
	clients := make([]otm.Client, 0, 2)
	for _, m := range s.config.COAP.OwnershipTransfer.Methods {
		var c otm.Client
		if m == OwnershipTransferManufacturerCertificate {
			c = s.getManufacturerCertificateClient()
		} else {
			c = s.getJustWorksClient()
		}
		clients = append(clients, c)
	}
	return clients
}

func (s *Service) getJustWorksClient() *justworks.Client {
	return justworks.NewClient(justworks.WithDialDTLS(s.DialOwnership))
}

func (s *Service) getManufacturerCertificateClient() *manufacturer.Client {
	return manufacturer.NewClient(s.config.COAP.OwnershipTransfer.Manufacturer.TLS.certificate, s.config.COAP.OwnershipTransfer.Manufacturer.TLS.caPool, manufacturer.WithDialDTLS(s.DialOwnership))
}

func (s *Service) GetOwnOptions() []core.OwnOption {
	return s.authenticationClient.GetOwnOptions()
}

// Serve starts a device provisioning service on the configured address in *Service.
func (s *Service) Serve() error {
	return s.serveWithHandlingSignal()
}

func (s *Service) serveWithHandlingSignal() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	services := make([]func() error, 0, 3)
	if s.udp4server != nil {
		services = append(services, func() error {
			return s.udp4server.Serve(s.udp4Listener)
		})
	}
	if s.udp6server != nil {
		services = append(services, func() error {
			return s.udp6server.Serve(s.udp6Listener)
		})
	}
	wg.Add(len(services))
	for _, serve := range services {
		go func(serve func() error) {
			defer wg.Done()
			err := serve()
			errCh <- err
		}(serve)
	}

	// Wait for a signal to shutdown
	<-s.done

	if s.udp4server != nil {
		s.udp4server.Stop()
	}
	if s.udp6server != nil {
		s.udp6server.Stop()
	}
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

// Close turn off server.
func (s *Service) Close() error {
	s.authenticationClient.Reset()
	close(s.done)
	return nil
}

func (s *Service) GetIdentityCSR(id string) ([]byte, error) {
	return s.authenticationClient.GetIdentityCSR(id)
}

func (s *Service) SetIdentityCertificate(chainPem []byte) error {
	return s.authenticationClient.SetIdentityCertificate(chainPem)
}

func (s *Service) GetIdentityCertificate() (tls.Certificate, error) {
	return s.authenticationClient.GetIdentityCertificate()
}

func (s *Service) GetDeviceAuthenticationMode() pb.GetConfigurationResponse_DeviceAuthenticationMode {
	switch s.config.COAP.TLS.Authentication {
	case AuthenticationX509:
		return pb.GetConfigurationResponse_X509
	case AuthenticationPreSharedKey:
		return pb.GetConfigurationResponse_PRE_SHARED_KEY
	}
	return pb.GetConfigurationResponse_PRE_SHARED_KEY
}

func (s *Service) IsInitialized() bool {
	return s.authenticationClient.IsInitialized()
}

func (s *Service) Reset() {
	s.authenticationClient.Reset()
}
