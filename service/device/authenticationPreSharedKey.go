package device

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/pkg/net/coap"
)

type authenticationPreSharedKey struct {
	config Config
}

var errPreSharedKeyAuthentication = fmt.Errorf("authentication method is set to %v", AuthenticationPreSharedKey)

func newAuthenticationPreSharedKey(config Config) *authenticationPreSharedKey {
	return &authenticationPreSharedKey{
		config: config,
	}
}

func (s *authenticationPreSharedKey) DialDTLS(ctx context.Context, addr string, _ *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	idBin, _ := s.config.COAP.TLS.PreSharedKey.subjectUUID.MarshalBinary()
	dtlsCfg := &dtls.Config{
		PSKIdentityHint: idBin,
		PSK: func(b []byte) ([]byte, error) {
			// iotivity-lite supports only 16-byte PSK
			return s.config.COAP.TLS.PreSharedKey.keyUUID[:16], nil
		},
		CipherSuites: []dtls.CipherSuiteID{dtls.TLS_ECDHE_PSK_WITH_AES_128_CBC_SHA256},
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, opts...)
}

func (s *authenticationPreSharedKey) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) GetOwnerID() (string, error) {
	return s.config.COAP.TLS.PreSharedKey.SubjectUUID, nil
}

func (s *authenticationPreSharedKey) GetOwnOptions() []core.OwnOption {
	return []core.OwnOption{core.WithPresharedKey(s.config.COAP.TLS.PreSharedKey.keyUUID[:16])}
}

func (s *authenticationPreSharedKey) GetCSR(id string) ([]byte, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) SetCertificate(chainPem []byte) error {
	return errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) GetCertificate() (tls.Certificate, error) {
	// we need to set empty certificate otherwise own device will failed
	return tls.Certificate{}, nil
}

func (s *authenticationPreSharedKey) GetCertificateAuthorities() ([]*x509.Certificate, error) {
	// we need to set empty certificates authorities otherwise own device will failed
	return nil, nil
}
