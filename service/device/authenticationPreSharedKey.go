package device

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/google/uuid"
	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v3/tcp"
	"github.com/plgd-dev/go-coap/v3/udp"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authenticationPreSharedKey struct {
	config *atomic.Pointer[Config]
}

var errPreSharedKeyAuthentication = status.Errorf(codes.Unimplemented, "authentication method is set to %v: not supported", AuthenticationPreSharedKey)

func newAuthenticationPreSharedKey(config Config) *authenticationPreSharedKey {
	return &authenticationPreSharedKey{
		config: atomic.NewPointer(&config),
	}
}

func (s *authenticationPreSharedKey) DialDTLS(ctx context.Context, addr string, _ *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	idBin, _ := s.config.Load().COAP.TLS.PreSharedKey.subjectUUID.MarshalBinary()
	dtlsCfg := &dtls.Config{
		PSKIdentityHint: idBin,
		PSK: func(b []byte) ([]byte, error) {
			// iotivity-lite supports only 16-byte PSK
			return s.config.Load().COAP.TLS.PreSharedKey.keyUUID[:16], nil
		},
		CipherSuites: []dtls.CipherSuiteID{dtls.TLS_ECDHE_PSK_WITH_AES_128_CBC_SHA256},
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, opts...)
}

func (s *authenticationPreSharedKey) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...tcp.Option) (*coap.ClientCloseHandler, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) GetOwnerID() (string, error) {
	return s.config.Load().COAP.TLS.PreSharedKey.SubjectUUID, nil
}

func (s *authenticationPreSharedKey) GetOwnOptions() []core.OwnOption {
	return []core.OwnOption{core.WithPresharedKey(s.config.Load().COAP.TLS.PreSharedKey.keyUUID[:16])}
}

func (s *authenticationPreSharedKey) GetIdentityCSR(id string) ([]byte, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) SetIdentityCertificate(owner string, chainPem []byte) error {
	return errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) GetIdentityCertificate() (tls.Certificate, error) {
	// we need to set empty certificate otherwise own device will failed
	return tls.Certificate{}, nil
}

func (s *authenticationPreSharedKey) GetCertificateAuthorities() ([]*x509.Certificate, error) {
	// we need to set empty certificates authorities otherwise own device will failed
	return nil, nil
}

func (s *authenticationPreSharedKey) IsInitialized() bool {
	cfg := s.config.Load()
	return cfg.COAP.TLS.PreSharedKey.keyUUID != uuid.Nil && cfg.COAP.TLS.PreSharedKey.subjectUUID != uuid.Nil
}

func (s *authenticationPreSharedKey) Reset() {
	// nothing to do
}

func (s *authenticationPreSharedKey) GetOwner() string {
	return s.config.Load().COAP.TLS.PreSharedKey.SubjectUUID
}
