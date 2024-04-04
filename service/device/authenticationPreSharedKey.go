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

	"github.com/google/uuid"
	"github.com/pion/dtls/v2"
	configDevice "github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v3/tcp"
	"github.com/plgd-dev/go-coap/v3/udp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authenticationPreSharedKey struct {
	getConfig func() configDevice.Config
}

var errPreSharedKeyAuthentication = status.Errorf(codes.Unimplemented, "authentication method is set to %v: not supported", configDevice.AuthenticationPreSharedKey)

func newAuthenticationPreSharedKey(getConfig func() configDevice.Config) *authenticationPreSharedKey {
	return &authenticationPreSharedKey{
		getConfig: getConfig,
	}
}

func (s *authenticationPreSharedKey) GetPreSharedKey() (uuid.UUID, string) {
	psk := s.getConfig().COAP.TLS.PreSharedKey
	return psk.Get()
}

func toKeyBin(key string) []byte {
	keyBin := []byte(key)
	if len(keyBin) < 16 {
		keyBin = append(keyBin, make([]byte, 16-len(keyBin))...)
	}
	return keyBin[:16]
}

func (s *authenticationPreSharedKey) DialDTLS(ctx context.Context, addr string, _ *dtls.Config, opts ...udp.Option) (*coap.ClientCloseHandler, error) {
	subjectUUID, key := s.GetPreSharedKey()
	if subjectUUID == uuid.Nil {
		return nil, status.Errorf(codes.Unauthenticated, "subjectId is empty")
	}
	if key == "" {
		return nil, status.Errorf(codes.Unauthenticated, "key is empty")
	}
	idBin, _ := subjectUUID.MarshalBinary()
	dtlsCfg := &dtls.Config{
		PSKIdentityHint: idBin,
		PSK: func([]byte) ([]byte, error) {
			// iotivity-lite supports only 16-byte PSK
			return toKeyBin(key), nil
		},
		CipherSuites: []dtls.CipherSuiteID{dtls.TLS_ECDHE_PSK_WITH_AES_128_CBC_SHA256},
	}
	return coap.DialUDPSecure(ctx, addr, dtlsCfg, opts...)
}

func (s *authenticationPreSharedKey) DialTLS(_ context.Context, _ string, _ *tls.Config, _ ...tcp.Option) (*coap.ClientCloseHandler, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) GetOwnerID() (string, error) {
	subjectUUID, _ := s.GetPreSharedKey()
	return subjectUUID.String(), nil
}

func (s *authenticationPreSharedKey) GetOwnOptions() ([]core.OwnOption, error) {
	_, key := s.GetPreSharedKey()
	if key == "" {
		return nil, status.Errorf(codes.Unauthenticated, "key is empty")
	}
	return []core.OwnOption{core.WithPresharedKey(toKeyBin(key))}, nil
}

func (s *authenticationPreSharedKey) GetIdentityCSR(_ string) ([]byte, error) {
	return nil, errPreSharedKeyAuthentication
}

func (s *authenticationPreSharedKey) SetIdentityCertificate(_ string, _ []byte) error {
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
	subjectUUID, key := s.GetPreSharedKey()
	return key != "" && subjectUUID != uuid.Nil
}

func (s *authenticationPreSharedKey) Reset() {
	// nothing to do
}

func (s *authenticationPreSharedKey) GetOwner() string {
	return s.getConfig().COAP.TLS.PreSharedKey.SubjectIDStr
}
