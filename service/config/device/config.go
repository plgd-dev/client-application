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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/go-coap/v3/net/blockwise"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/security/certManager/client"
	pkgStrings "github.com/plgd-dev/hub/v2/pkg/strings"
	"github.com/plgd-dev/kit/v2/security"
)

type Config struct {
	COAP CoapConfig `yaml:"coap" json:"coap"`
}

func (c *Config) Validate() error {
	if err := c.COAP.Validate(); err != nil {
		return fmt.Errorf("coap.%w", err)
	}
	return nil
}

type ManufacturerTLSConfig struct {
	CAPool      interface{}         `yaml:"caPool" json:"caPool" description:"file path to the root certificates in PEM format"`
	KeyFile     string              `yaml:"keyFile" json:"keyFile" description:"file name of private key in PEM format"`
	CertFile    string              `yaml:"certFile" json:"certFile" description:"file name of certificate in PEM format"`
	certificate tls.Certificate     `yaml:"-"`
	caPool      []*x509.Certificate `yaml:"-"`
}

func (c *ManufacturerTLSConfig) GetCAPool() []*x509.Certificate {
	return c.caPool
}

func (c *ManufacturerTLSConfig) GetCertificate() tls.Certificate {
	return c.certificate
}

func (c *ManufacturerTLSConfig) Validate() error {
	caPoolArray, ok := pkgStrings.ToStringArray(c.CAPool)
	if !ok {
		return fmt.Errorf("caPool('%v')", c.CAPool)
	}
	var caPool []*x509.Certificate
	for idx, ca := range caPoolArray {
		certs, err := security.LoadX509(ca)
		if err != nil {
			return fmt.Errorf("caPool[%d]('%v') - %w", idx, ca, err)
		}
		caPool = append(caPool, certs...)
	}
	certificate, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return fmt.Errorf("certFile('%v'), keyFile('%v') - %w", c.CertFile, c.KeyFile, err)
	}
	c.caPool = caPool
	c.certificate = certificate
	return nil
}

func (c *ManufacturerTLSConfig) ToCertMangerConfig() client.Config {
	return client.Config{
		CAPool:          c.CAPool,
		KeyFile:         c.KeyFile,
		CertFile:        c.CertFile,
		UseSystemCAPool: false,
	}
}

type ManufacturerConfig struct {
	TLS ManufacturerTLSConfig `yaml:"tls" json:"tls"`
}

func (c *ManufacturerConfig) Validate() error {
	if err := c.TLS.Validate(); err != nil {
		return fmt.Errorf("tls.%w", err)
	}
	return nil
}

type OwnershipTransferMethod string

const (
	OwnershipTransferJustWorks               OwnershipTransferMethod = "justWorks"
	OwnershipTransferManufacturerCertificate OwnershipTransferMethod = "manufacturerCertificate"
)

var validOwnershipTransfers = map[OwnershipTransferMethod]bool{
	OwnershipTransferJustWorks:               true,
	OwnershipTransferManufacturerCertificate: true,
}

type OwnershipTransferConfig struct {
	Methods      []OwnershipTransferMethod `yaml:"methods" json:"methods"`
	Manufacturer ManufacturerConfig        `yaml:"manufacturerCertificate" json:"manufacturerCertificate"`
}

func (c *OwnershipTransferConfig) Validate() error {
	containsManufacturerCertificate := false
	if len(c.Methods) == 0 {
		return fmt.Errorf("methods('%v') - is empty", c.Methods)
	}
	for idx, method := range c.Methods {
		if ok := validOwnershipTransfers[method]; !ok {
			return fmt.Errorf("methods[%v]('%v') - supports only '%v,%v'", idx, method, OwnershipTransferJustWorks, OwnershipTransferManufacturerCertificate)
		}
		if method == OwnershipTransferManufacturerCertificate {
			containsManufacturerCertificate = true
		}
	}
	if containsManufacturerCertificate {
		if err := c.Manufacturer.Validate(); err != nil {
			return fmt.Errorf("manufacturerCertificate.%w", err)
		}
	}

	return nil
}

type CoapConfig struct {
	MaxMessageSize    uint32                  `yaml:"maxMessageSize" json:"maxMessageSize"`
	InactivityMonitor InactivityMonitor       `yaml:"inactivityMonitor" json:"inactivityMonitor"`
	BlockwiseTransfer BlockwiseTransferConfig `yaml:"blockwiseTransfer" json:"blockwiseTransfer"`
	OwnershipTransfer OwnershipTransferConfig `yaml:"ownershipTransfer" json:"ownershipTransfer"`
	TLS               TLSConfig               `yaml:"tls" json:"tls"`
}

func (c *CoapConfig) Validate() error {
	if c.MaxMessageSize <= 64 {
		return fmt.Errorf("maxMessageSize('%v')", c.MaxMessageSize)
	}
	if err := c.InactivityMonitor.Validate(); err != nil {
		return fmt.Errorf("inactivityMonitor.%w", err)
	}
	if err := c.BlockwiseTransfer.Validate(); err != nil {
		return fmt.Errorf("blockwiseTransfer.%w", err)
	}
	if err := c.OwnershipTransfer.Validate(); err != nil {
		return fmt.Errorf("ownershipTransfer.%w", err)
	}
	if err := c.TLS.Validate(); err != nil {
		return fmt.Errorf("tls.%w", err)
	}
	return nil
}

type PreSharedKeyConfig struct {
	SubjectIDStr string    `yaml:"subjectId" json:"subjectId"`
	subjectID    uuid.UUID `yaml:"-"`
	Key          string    `yaml:"key" json:"key"`
}

func (c *PreSharedKeyConfig) Get() (uuid.UUID, string) {
	return c.subjectID, c.Key
}

func (c *PreSharedKeyConfig) Validate() error {
	var err error
	if c.Key == "" && c.SubjectIDStr == "" {
		c.subjectID = uuid.Nil
		return nil
	}
	if c.Key == "" {
		return fmt.Errorf("key('%v') - is empty", c.Key)
	}
	if c.subjectID, err = uuid.Parse(events.OwnerToUUID(c.SubjectIDStr)); err != nil || c.subjectID == uuid.Nil {
		return fmt.Errorf("subjectUUID('%v') - %w", c.SubjectIDStr, err)
	}
	return nil
}

type Authentication string

const (
	AuthenticationPreSharedKey Authentication = "preSharedKey"
	AuthenticationX509         Authentication = "x509"
)

type TLSConfig struct {
	Authentication Authentication     `yaml:"authentication" json:"authentication"`
	PreSharedKey   PreSharedKeyConfig `yaml:"preSharedKey" json:"preSharedKey"`
}

func (c *TLSConfig) Validate() error {
	switch c.Authentication {
	case AuthenticationX509:
	case AuthenticationPreSharedKey:
		if err := c.PreSharedKey.Validate(); err != nil {
			return fmt.Errorf("preSharedKey.%w", err)
		}
	default:
		return fmt.Errorf("authentication('%v') - supports only '%v,%v'", c.Authentication, AuthenticationPreSharedKey, AuthenticationX509)
	}
	return nil
}

type InactivityMonitor struct {
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

func (c *InactivityMonitor) Validate() error {
	if c.Timeout < time.Second {
		return fmt.Errorf("timeout('%v')", c.Timeout)
	}
	return nil
}

type BlockwiseTransferConfig struct {
	Enabled bool          `yaml:"enabled" json:"enabled"`
	SZXStr  string        `yaml:"blockSize" json:"blockSize"`
	szx     blockwise.SZX `yaml:"-" json:"-"`
}

func (c *BlockwiseTransferConfig) GetSZX() blockwise.SZX {
	return c.szx
}

func (c *BlockwiseTransferConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	switch strings.ToLower(c.SZXStr) {
	case "16":
		c.szx = blockwise.SZX16
	case "32":
		c.szx = blockwise.SZX32
	case "64":
		c.szx = blockwise.SZX64
	case "128":
		c.szx = blockwise.SZX128
	case "256":
		c.szx = blockwise.SZX256
	case "512":
		c.szx = blockwise.SZX512
	case "1024":
		c.szx = blockwise.SZX1024
	case "bert":
		c.szx = blockwise.SZXBERT
	default:
		return fmt.Errorf("blockSize('%v')", c.SZXStr)
	}
	return nil
}

var defaultConfig = Config{
	COAP: CoapConfig{
		MaxMessageSize: 256 * 1024,
		InactivityMonitor: InactivityMonitor{
			Timeout: time.Second * 10,
		},
		BlockwiseTransfer: BlockwiseTransferConfig{
			Enabled: true,
			SZXStr:  "1024",
		},
		TLS: TLSConfig{
			Authentication: AuthenticationPreSharedKey,
			PreSharedKey: PreSharedKeyConfig{
				SubjectIDStr: uuid.NewString(),
				Key:          uuid.NewString(),
			},
		},
		OwnershipTransfer: OwnershipTransferConfig{
			Methods: []OwnershipTransferMethod{OwnershipTransferJustWorks},
		},
	},
}

func DefaultConfig() Config {
	return defaultConfig
}
