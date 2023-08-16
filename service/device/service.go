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
	configDevice "github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/client/core/otm"
	justworks "github.com/plgd-dev/device/v2/client/core/otm/just-works"
	"github.com/plgd-dev/device/v2/client/core/otm/manufacturer"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/device/v2/schema"
	coapNet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/tcp"
	tcpClient "github.com/plgd-dev/go-coap/v3/tcp/client"
	"github.com/plgd-dev/go-coap/v3/udp"
	udpClient "github.com/plgd-dev/go-coap/v3/udp/client"
	udpServer "github.com/plgd-dev/go-coap/v3/udp/server"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

type AuthenticationClient interface {
	DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error)
	DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...tcp.Option) (*coap.ClientCloseHandler, error)
	GetOwnerID() (string, error)
	GetOwner() string
	GetOwnOptions() ([]core.OwnOption, error)

	GetIdentityCSR(id string) ([]byte, error)
	SetIdentityCertificate(owner string, chainPem []byte) error
	GetIdentityCertificate() (tls.Certificate, error)
	GetCertificateAuthorities() ([]*x509.Certificate, error)
	IsInitialized() bool
	Reset()
}

type Service struct {
	getConfig            func() configDevice.Config
	logger               log.Logger
	udp4server           *udpServer.Server
	udp6server           *udpServer.Server
	udp4Listener         *coapNet.UDPConn
	udp6Listener         *coapNet.UDPConn
	done                 chan struct{}
	authenticationClient AuthenticationClient
}

var closeHandlerKey = "close-handler"

func errClosingConnection(debugf func(fmt string, a ...any), scheme schema.Scheme, remoteAddr net.Addr) {
	debugf("closing connection %v://%v for inactivity", scheme, remoteAddr)
}

// New creates new GRPC service
func New(ctx context.Context, getConfig func() configDevice.Config, logger log.Logger) (*Service, error) {
	config := getConfig()
	var authenticationClient AuthenticationClient
	switch config.COAP.TLS.Authentication {
	case configDevice.AuthenticationPreSharedKey:
		authenticationClient = newAuthenticationPreSharedKey(getConfig)
	case configDevice.AuthenticationX509:
		authenticationClient = newAuthenticationX509(config)
	case configDevice.AuthenticationUninitialized:
		return nil, fmt.Errorf("device is not initialized")
	}

	opts := []udpServer.Option{
		options.WithContext(ctx),
		options.WithOnNewConn(func(cc *udpClient.Conn) {
			closeHander := coap.NewOnCloseHandler()
			cc.AddOnClose(func() {
				closeHander.OnClose(nil)
			})
			cc.SetContextValue(&closeHandlerKey, closeHander)
		}),
		options.WithMaxMessageSize(config.COAP.MaxMessageSize),
		options.WithBlockwise(config.COAP.BlockwiseTransfer.Enabled, config.COAP.BlockwiseTransfer.GetSZX(), config.COAP.InactivityMonitor.Timeout),
		options.WithErrors(errFunc),
		options.WithInactivityMonitor(config.COAP.InactivityMonitor.Timeout, func(cc *udpClient.Conn) {
			errClosingConnection(log.Debugf, "coap", cc.RemoteAddr())
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
		getConfig:            getConfig,
		logger:               logger,
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

func (s *Service) errFunc(err error) {
	s.logger.Debug(err)
}

func (s *Service) getDialUDPOptions(secure bool) []udp.Option {
	config := s.getConfig()
	dialOpts := []udp.Option{
		options.WithInactivityMonitor(config.COAP.InactivityMonitor.Timeout, func(cc *udpClient.Conn) {
			scheme := schema.UDPScheme
			if secure {
				scheme = schema.TCPSecureScheme
			}
			errClosingConnection(log.Debugf, scheme, cc.RemoteAddr())
			_ = cc.Close()
		}),
		options.WithMaxMessageSize(config.COAP.MaxMessageSize),
		options.WithBlockwise(config.COAP.BlockwiseTransfer.Enabled, config.COAP.BlockwiseTransfer.GetSZX(), config.COAP.InactivityMonitor.Timeout),
		options.WithErrors(s.errFunc),
	}
	return dialOpts
}

func (s *Service) DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	return s.authenticationClient.DialDTLS(ctx, addr, dtlsCfg, append(s.getDialUDPOptions(true), opts...)...)
}

func (s *Service) DialOwnership(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, append(s.getDialUDPOptions(true), opts...)...)
}

type UDPClientConn struct {
	*udpClient.Conn
}

func (c *UDPClientConn) Context() context.Context {
	cc := c.Conn
	// we need to check if connection will be closed for inactivity
	cc.CheckExpirations(time.Now().Add(50 * time.Millisecond))
	if cc.Context().Err() == nil {
		// move inactivity timeout to future
		cc.InactivityMonitor().Notify()
	}
	return cc.Context()
}

func (s *Service) DialUDP(ctx context.Context, addr string, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	var cc *udpClient.Conn
	if udpAddr.IP.To4() != nil {
		cc, err = s.udp4server.NewConn(udpAddr)
	} else {
		cc, err = s.udp6server.NewConn(udpAddr)
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
	return coap.NewClientCloseHandler(&UDPClientConn{Conn: cc}, h), nil
}

func (s *Service) getDialTCPOptions(secure bool) []tcp.Option {
	config := s.getConfig()
	dialOpts := []tcp.Option{
		options.WithInactivityMonitor(config.COAP.InactivityMonitor.Timeout, func(cc *tcpClient.Conn) {
			scheme := schema.TCPScheme
			if secure {
				scheme = schema.TCPSecureScheme
			}
			errClosingConnection(log.Debugf, scheme, cc.RemoteAddr())
			_ = cc.Close()
		}),
		options.WithMaxMessageSize(config.COAP.MaxMessageSize),
		options.WithBlockwise(config.COAP.BlockwiseTransfer.Enabled, config.COAP.BlockwiseTransfer.GetSZX(), config.COAP.InactivityMonitor.Timeout),
		options.WithErrors(s.errFunc),
	}
	return dialOpts
}

func (s *Service) DialTCP(ctx context.Context, addr string, opts ...tcp.Option) (*coap.ClientCloseHandler, error) {
	return coap.DialTCP(ctx, addr, append(s.getDialTCPOptions(false), opts...)...)
}

func (s *Service) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...tcp.Option) (*coap.ClientCloseHandler, error) {
	return s.authenticationClient.DialTLS(ctx, addr, tlsCfg, append(s.getDialTCPOptions(true), opts...)...)
}

func (s *Service) DeviceLogger() core.Logger {
	return s.logger.DTLSLoggerFactory().NewLogger("client-application/device")
}

func (s *Service) GetDeviceConfiguration() core.DeviceConfiguration {
	return core.DeviceConfiguration{
		DialDTLS: s.DialDTLS,
		DialUDP:  s.DialUDP,
		DialTCP:  s.DialTCP,
		DialTLS:  s.DialTLS,
		Logger:   s.DeviceLogger(),
		TLSConfig: &core.TLSConfig{
			GetCertificate:            s.authenticationClient.GetIdentityCertificate,
			GetCertificateAuthorities: s.authenticationClient.GetCertificateAuthorities,
		},
		GetOwnerID: s.authenticationClient.GetOwnerID,
	}
}

func (s *Service) GetOwnershipClients() []otm.Client {
	clients := make([]otm.Client, 0, 2)
	config := s.getConfig()
	for _, m := range config.COAP.OwnershipTransfer.Methods {
		var c otm.Client
		if m == configDevice.OwnershipTransferManufacturerCertificate {
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
	config := s.getConfig()
	return manufacturer.NewClient(config.COAP.OwnershipTransfer.Manufacturer.TLS.GetCertificate(), config.COAP.OwnershipTransfer.Manufacturer.TLS.GetCAPool(), manufacturer.WithDialDTLS(s.DialOwnership))
}

func (s *Service) GetOwnOptions() ([]core.OwnOption, error) {
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

func (s *Service) SetIdentityCertificate(owner string, chainPem []byte) error {
	return s.authenticationClient.SetIdentityCertificate(owner, chainPem)
}

func (s *Service) GetIdentityCertificate() (tls.Certificate, error) {
	return s.authenticationClient.GetIdentityCertificate()
}

func (s *Service) GetDeviceAuthenticationMode() pb.GetConfigurationResponse_DeviceAuthenticationMode {
	config := s.getConfig()
	switch config.COAP.TLS.Authentication {
	case configDevice.AuthenticationX509:
		return pb.GetConfigurationResponse_X509
	case configDevice.AuthenticationPreSharedKey:
		return pb.GetConfigurationResponse_PRE_SHARED_KEY
	case configDevice.AuthenticationUninitialized:
		return pb.GetConfigurationResponse_UNINITIALIZED
	}
	return pb.GetConfigurationResponse_UNINITIALIZED
}

func (s *Service) IsInitialized() bool {
	return s.authenticationClient.IsInitialized()
}

func (s *Service) Reset() {
	s.authenticationClient.Reset()
}

func (s *Service) GetOwner() string {
	return s.authenticationClient.GetOwner()
}
