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
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v3/tcp"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/kit/v2/security/generateCertificate"
	"go.uber.org/atomic"
)

type authenticationX509 struct {
	config      Config
	privateKey  atomic.Pointer[ecdsa.PrivateKey]
	certificate atomic.Pointer[tls.Certificate]
	owner       atomic.String
}

func newAuthenticationX509(config Config) *authenticationX509 {
	return &authenticationX509{
		config: config,
	}
}

func (s *authenticationX509) getTLSCertificate() (*tls.Certificate, error) {
	crt := s.certificate.Load()
	if crt == nil || crt.Leaf == nil {
		return nil, fmt.Errorf("certificate hasn't been set")
	}
	if crt.Leaf.NotAfter.Before(time.Now()) {
		return nil, fmt.Errorf("certificate is not valid after %v", crt.Leaf.NotAfter)
	}
	if crt.Leaf.NotBefore.After(time.Now()) {
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

func (s *authenticationX509) DialDTLS(ctx context.Context, addr string, dtlsCfg *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	crt, clientCAs, err := s.getClientCerts()
	if err != nil {
		return nil, err
	}
	dtlsCfg.Certificates = []tls.Certificate{*crt}
	dtlsCfg.ClientCAs = clientCAs
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, opts...)
}

func (s *authenticationX509) DialTLS(ctx context.Context, addr string, tlsCfg *tls.Config, opts ...tcp.Option) (*coap.ClientCloseHandler, error) {
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

func (s *authenticationX509) GetIdentityCSR(id string) ([]byte, error) {
	var privateKey *ecdsa.PrivateKey
	for {
		privateKey = s.privateKey.Load()
		if privateKey != nil {
			break
		}
		var err error
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("cannot generate private key: %w", err)
		}
		if s.privateKey.CompareAndSwap(nil, privateKey) {
			break
		}
	}
	return generateCertificate.GenerateIdentityCSR(generateCertificate.Configuration{}, id, privateKey)
}

func encodePrivateKeyToPem(k *ecdsa.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b}), nil
}

func (s *authenticationX509) setCertificate(crt tls.Certificate) (bool, error) {
	oldCrt := s.certificate.Load()
	if oldCrt == nil {
		if s.certificate.CompareAndSwap(oldCrt, &crt) {
			return true, nil
		}
		return false, nil
	}
	if crt.Leaf.Subject.CommonName != oldCrt.Leaf.Subject.CommonName {
		return false, fmt.Errorf("common name of certificate(%v) is not equal with previous one(%v)", crt.Leaf.Subject.CommonName, oldCrt.Leaf.Subject.CommonName)
	}
	if s.certificate.CompareAndSwap(oldCrt, &crt) {
		return true, nil
	}
	return false, nil
}

func (s *authenticationX509) updateCertificate(crt tls.Certificate) error {
	if len(crt.Certificate) > 0 && crt.Leaf == nil {
		var err error
		crt.Leaf, err = x509.ParseCertificate(crt.Certificate[0])
		if err != nil {
			return fmt.Errorf("cannot parse leaf certificate: %w", err)
		}
	}
	for {
		ok, err := s.setCertificate(crt)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
}

func (s *authenticationX509) SetIdentityCertificate(owner string, chainPem []byte) error {
	privateKey := s.privateKey.Load()
	if privateKey == nil {
		return fmt.Errorf("private key is not set")
	}
	keyPem, err := encodePrivateKeyToPem(privateKey)
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
	if err := s.updateCertificate(crt); err != nil {
		return fmt.Errorf("cannot update certificate: %w", err)
	}
	s.owner.Store(owner)
	return nil
}

func (s *authenticationX509) GetIdentityCertificate() (tls.Certificate, error) {
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

func (s *authenticationX509) IsInitialized() bool {
	return s.certificate.Load() != nil
}

func (s *authenticationX509) Reset() {
	s.certificate.Store((*tls.Certificate)(nil))
	s.privateKey.Store((*ecdsa.PrivateKey)(nil))
	s.owner.Store("")
}

func (s *authenticationX509) GetOwner() string {
	return s.owner.Load()
}
