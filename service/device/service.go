package device

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/client/core"
	justworks "github.com/plgd-dev/device/client/core/otm/just-works"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/message/status"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	config         Config
	logger         log.Logger
	tracerProvider trace.TracerProvider
}

// New creates new GRPC service
func New(ctx context.Context, serviceName string, config Config, logger log.Logger, tracerProvider trace.TracerProvider) (*Service, error) {
	return &Service{
		config:         config,
		logger:         logger,
		tracerProvider: tracerProvider,
	}, nil
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
	dialOpts := []coap.DialOptionFunc{
		coap.WithInactivityMonitor(s.config.COAP.InactivityMonitor.Timeout),
		coap.WithMaxMessageSize(s.config.COAP.MaxMessageSize),
		coap.WithBlockwise(s.config.COAP.BlockwiseTransfer.Enabled, s.config.COAP.BlockwiseTransfer.szx, s.config.COAP.InactivityMonitor.Timeout),
		coap.WithErrors(s.ErrFunc),
	}
	return coap.DialUDP(ctx, addr, append(dialOpts, opts...)...)
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
