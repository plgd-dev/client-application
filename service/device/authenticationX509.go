package device

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/kit/v2/security/generateCertificate"
	"go.uber.org/atomic"
)

type authenticationX509 struct {
	config      Config
	privateKey  *ecdsa.PrivateKey
	certificate atomic.Pointer[tls.Certificate]
}

func newAuthenticationX509(config Config) (*authenticationX509, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("cannot generate private key: %w", err)
	}
	return &authenticationX509{
		config:     config,
		privateKey: privateKey,
	}, nil
}

func (s *authenticationX509) getTLSCertificate() (*tls.Certificate, error) {
	crt := s.certificate.Load()
	if crt == nil || crt.Leaf == nil {
		return nil, fmt.Errorf("certificate hasn't been set")
	}
	if crt.Leaf.NotAfter.After(time.Now()) {
		return nil, fmt.Errorf("certificate is not valid after %v", crt.Leaf.NotAfter)
	}
	if crt.Leaf.NotBefore.Before(time.Now()) {
		return nil, fmt.Errorf("certificate is not valid before %v", crt.Leaf.NotBefore)
	}
	return crt, nil
}

func getRootCAFromChain(chain [][]byte) (*x509.Certificate, error) {
	rootCA, err := x509.ParseCertificate(chain[len(chain)-1])
	if err != nil {
		return nil, err
	}
	if !rootCA.IsCA || rootCA.Issuer.CommonName != rootCA.Subject.CommonName {
		return nil, fmt.Errorf("invalid root certificate")
	}
	return rootCA, nil
}

func (s *authenticationX509) getClientCerts() (*tls.Certificate, *x509.CertPool, error) {
	crt, err := s.getTLSCertificate()
	if err != nil {
		return nil, nil, err
	}
	rootCA, err := getRootCAFromChain(crt.Certificate)
	if err != nil {
		return nil, nil, err
	}
	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(rootCA)
	return crt, clientCAs, nil
}

func (s *authenticationX509) DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	crt, clientCAs, err := s.getClientCerts()
	if err != nil {
		return nil, err
	}
	dtlsCfg.Certificates = []tls.Certificate{*crt}
	dtlsCfg.ClientCAs = clientCAs
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, opts...)
}

func (s *authenticationX509) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...coap.DialOptionFunc) (*coap.ClientCloseHandler, error) {
	crt, clientCAs, err := s.getClientCerts()
	if err != nil {
		return nil, err
	}
	tlsCfg.Certificates = []tls.Certificate{*crt}
	tlsCfg.ClientCAs = clientCAs
	return coap.DialTCPSecure(ctx, addr, tlsCfg, opts...)
}

func (s *authenticationX509) GetOwnerID() (string, error) {
	crt, err := s.getTLSCertificate()
	if err != nil {
		return "", err
	}
	return coap.GetDeviceIDFromIdentityCertificate(crt.Leaf)
}

func (s *authenticationX509) GetOwnOptions() []core.OwnOption {
	return nil
}

func (s *authenticationX509) GetCSR(id string) ([]byte, error) {
	return generateCertificate.GenerateIdentityCSR(generateCertificate.Configuration{}, id, s.privateKey)
}

func encodePrivateKeyToPem(k *ecdsa.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b}), nil
}

func (s *authenticationX509) SetCertificate(chainPem []byte) error {
	keyPem, err := encodePrivateKeyToPem(s.privateKey)
	if err != nil {
		return fmt.Errorf("cannot marshal private key: %w", err)
	}

	crt, err := tls.X509KeyPair(chainPem, keyPem)
	if err != nil {
		return fmt.Errorf("cannot create certificate: %w", err)
	}
	if _, err := getRootCAFromChain(crt.Certificate); err != nil {
		return fmt.Errorf("invalid root certificate: %w", err)
	}
	s.certificate.Store(&crt)
	return nil
}

func (s *authenticationX509) GetCertificate() (tls.Certificate, error) {
	crt, err := s.getTLSCertificate()
	if err != nil {
		return tls.Certificate{}, err
	}
	return *crt, nil
}

func (s *authenticationX509) GetCertificateAuthorities() ([]*x509.Certificate, error) {
	crt, err := s.getTLSCertificate()
	if err != nil {
		return nil, err
	}
	rootCA, err := getRootCAFromChain(crt.Certificate)
	if err != nil {
		return nil, err
	}
	return []*x509.Certificate{rootCA}, nil
}
