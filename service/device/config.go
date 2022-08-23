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
	"github.com/plgd-dev/go-coap/v2/net/blockwise"
	"github.com/plgd-dev/hub/v2/pkg/security/certManager/client"
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
	CAPool      string              `yaml:"caPool" json:"caPool" description:"file path to the root certificates in PEM format"`
	KeyFile     string              `yaml:"keyFile" json:"keyFile" description:"file name of private key in PEM format"`
	CertFile    string              `yaml:"certFile" json:"certFile" description:"file name of certificate in PEM format"`
	certificate tls.Certificate     `yaml:"-"`
	caPool      []*x509.Certificate `yaml:"-"`
}

func (c *ManufacturerTLSConfig) Validate() error {
	caPool, err := security.LoadX509(c.CAPool)
	if err != nil {
		return fmt.Errorf("caPool('%v') - %w", c.CAPool, err)
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
	OwnershipTransferJustWorks    OwnershipTransferMethod = "justWorks"
	OwnershipTransferManufacturer OwnershipTransferMethod = "manufacturer"
)

var validOwnershipTransfers = map[OwnershipTransferMethod]bool{
	OwnershipTransferJustWorks:    true,
	OwnershipTransferManufacturer: true,
}

type OwnershipTransferConfig struct {
	Method       OwnershipTransferMethod `yaml:"method" json:"method"`
	Manufacturer ManufacturerConfig      `yaml:"manufacturer" json:"manufacturer"`
}

func (c *OwnershipTransferConfig) Validate() error {
	if ok := validOwnershipTransfers[c.Method]; !ok {
		return fmt.Errorf("method('%v') - supports only '%v,%v'", c.Method, OwnershipTransferJustWorks, OwnershipTransferManufacturer)
	}
	switch c.Method {
	case OwnershipTransferJustWorks:
		return nil
	case OwnershipTransferManufacturer:
		if err := c.Manufacturer.Validate(); err != nil {
			return fmt.Errorf("manufacturer.%w", err)
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

type TLSConfig struct {
	SubjectUUID      string `yaml:"subjectUuid" json:"subjectUuid"`
	subjectUUID      uuid.UUID
	PreSharedKeyUUID string `yaml:"preSharedKeyUuid" json:"preSharedKeyUuid"`
	preSharedKeyUUID uuid.UUID
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
	SZX     string        `yaml:"blockSize" json:"blockSize"`
	szx     blockwise.SZX `yaml:"-" json:"-"`
}

func (c *BlockwiseTransferConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	switch strings.ToLower(c.SZX) {
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
		return fmt.Errorf("blockSize('%v')", c.SZX)
	}
	return nil
}

func (c *TLSConfig) Validate() error {
	var err error
	if c.preSharedKeyUUID, err = uuid.Parse(c.PreSharedKeyUUID); err != nil {
		return fmt.Errorf("preSharedKeyUUID('%v') - %w", c.PreSharedKeyUUID, err)
	}
	if c.subjectUUID, err = uuid.Parse(c.SubjectUUID); err != nil {
		return fmt.Errorf("subjectUUID('%v') - %w", c.SubjectUUID, err)
	}
	return nil
}
