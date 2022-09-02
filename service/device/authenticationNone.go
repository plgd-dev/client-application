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

type authenticationNone struct{}

var errNoneAuthentication = fmt.Errorf("authentication method is set to %v", AuthenticationNone)

func newAuthenticationNone() *authenticationNone {
	return &authenticationNone{}
}

func (s *authenticationNone) DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	return nil, errNoneAuthentication
}

func (s *authenticationNone) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	return nil, errNoneAuthentication
}

func (s *authenticationNone) GetOwnerID() (string, error) {
	return "", errNoneAuthentication
}

func (s *authenticationNone) GetOwnOptions() []core.OwnOption {
	return nil
}

func (s *authenticationNone) GetCSR(id string) ([]byte, error) {
	return nil, errNoneAuthentication
}

func (s *authenticationNone) SetCertificate(chainPem []byte) error {
	return errNoneAuthentication
}

func (s *authenticationNone) GetCertificate() (tls.Certificate, error) {
	return tls.Certificate{}, errNoneAuthentication
}

func (s *authenticationNone) GetCertificateAuthorities() ([]*x509.Certificate, error) {
	return nil, errNoneAuthentication
}
