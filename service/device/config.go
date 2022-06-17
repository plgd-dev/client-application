//***************************************************************************
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
//*************************************************************************

package device

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/go-coap/v2/net/blockwise"
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

type CoapConfig struct {
	MaxMessageSize    uint32                  `yaml:"maxMessageSize" json:"maxMessageSize"`
	InactivityMonitor InactivityMonitor       `yaml:"inactivityMonitor" json:"inactivityMonitor"`
	BlockwiseTransfer BlockwiseTransferConfig `yaml:"blockwiseTransfer" json:"blockwiseTransfer"`
	TLS               TLSConfig               `yaml:"tls" json:"tls"`
}

func (c *CoapConfig) Validate() error {
	if c.MaxMessageSize <= 64 {
		return fmt.Errorf("maxMessageSize('%v')", c.MaxMessageSize)
	}
	if err := c.InactivityMonitor.Validate(); err != nil {
		return fmt.Errorf("keepAlive.%w", err)
	}
	if err := c.BlockwiseTransfer.Validate(); err != nil {
		return fmt.Errorf("blockwiseTransfer.%w", err)
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
