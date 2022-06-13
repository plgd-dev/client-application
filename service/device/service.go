package device

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"sync"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/client/core"
	justworks "github.com/plgd-dev/device/client/core/otm/just-works"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/message/status"
	coapNet "github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/net/monitor/inactivity"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	config         Config
	logger         log.Logger
	tracerProvider trace.TracerProvider
	udp4server     *udp.Server
	udp6server     *udp.Server
	udp4Listener   *coapNet.UDPConn
	udp6Listener   *coapNet.UDPConn
	done           chan struct{}
}

var closeHandlerKey = "close-handler"

// New creates new GRPC service
func New(ctx context.Context, serviceName string, config Config, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	opts := []udp.ServerOption{
		udp.WithContext(ctx),
		udp.WithOnNewClientConn(func(cc *client.ClientConn) {
			closeHander := coap.NewOnCloseHandler()
			cc.AddOnClose(func() {
				closeHander.OnClose(nil)
			})
			cc.SetContextValue(&closeHandlerKey, coap.NewOnCloseHandler())
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
		config:         config,
		logger:         logger,
		tracerProvider: tracerProvider,
		udp4server:     udp4server,
		udp6server:     udp6server,
		udp4Listener:   udp4Listener,
		udp6Listener:   udp6Listener,
		done:           make(chan struct{}),
	}, nil
}

func errFunc(err error) {
	log.Debug(err)
}

func (s *Service) DialDTLS(ctx context.Context, addr string, _ *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	idBin, _ := s.config.COAP.TLS.subjectUUID.MarshalBinary()
	dtlsCfg := &dtls.Config{
		PSKIdentityHint: idBin,
		PSK: func(b []byte) ([]byte, error) {
			// iotivity-lite supports only 16-byte PSK
			return s.config.COAP.TLS.preSharedKeyUUID[:16], nil
		},
		CipherSuites: []dtls.CipherSuiteID{dtls.TLS_ECDHE_PSK_WITH_AES_128_CBC_SHA256},
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, append(dialOpts, opts...)...)
}

func (s *Service) DialJustWorks(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, append(dialOpts, opts...)...)
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
	return coap.NewClientCloseHandler(cc.Client(), h), nil
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
	return nil, status.Errorf(&message.Message{
		Code: codes.NotImplemented,
	}, "not implemented")
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
			GetCertificate: func() (tls.Certificate, error) {
				return tls.Certificate{}, nil
			},
			GetCertificateAuthorities: func() ([]*x509.Certificate, error) {
				return nil, nil
			},
		},
		GetOwnerID: func() (string, error) {
			return s.config.COAP.TLS.SubjectUUID, nil
		},
	}
}

func (s *Service) GetJustWorksClient() *justworks.Client {
	return justworks.NewClient(justworks.WithDialDTLS(s.DialJustWorks))
}

func (s *Service) GetOwnOptions() []core.OwnOption {
	return []core.OwnOption{core.WithPresharedKey(s.config.COAP.TLS.preSharedKeyUUID[:16])}
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

// Shutdown turn off server.
func (s *Service) Close() {
	close(s.done)
}
