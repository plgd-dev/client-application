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

package server

import (
	"fmt"
	"time"

	"github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
	certManager "github.com/plgd-dev/hub/v2/pkg/security/certManager/server"
)

type TLSConfig struct {
	Enabled            bool `yaml:"enabled" json:"enabled"`
	certManager.Config `yaml:",inline"`
}

func (c *TLSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	return c.Config.Validate()
}

type UserAgentConfig struct {
	CertificateAuthorityAddress string        `yaml:"certificateAuthorityAddress" json:"certificateAuthorityAddress"`
	CSRExpiration               time.Duration `yaml:"csrExpiration" json:"csrExpiration"`
}

func (c *UserAgentConfig) Validate() error {
	if c.CertificateAuthorityAddress == "" {
		return fmt.Errorf("certificateAuthorityAddress('%v')", c.CertificateAuthorityAddress)
	}
	if c.CSRExpiration < 0 {
		return fmt.Errorf("csrExpiration('%v')", c.CSRExpiration)
	}
	return nil
}

type MediatorMode string

const (
	MediatorMode_None      MediatorMode = "none"
	MediatorMode_UserAgent MediatorMode = "userAgent"
)

// MediatorConfig represents web configuration for user interface exposed via getOAuthConfiguration handler
type MediatorConfig struct {
	Mode            MediatorMode    `yaml:"mode" json:"mode"`
	UserAgentConfig UserAgentConfig `yaml:"userAgent" json:"userAgent"`
}

func (c *MediatorConfig) Validate() error {
	switch c.Mode {
	case MediatorMode_None, "":
		c.Mode = MediatorMode_None
	case MediatorMode_UserAgent:
		if err := c.UserAgentConfig.Validate(); err != nil {
			return fmt.Errorf("userAgent.%w", err)
		}
	}
	return nil
}

type AuthorizationConfig struct {
	ClientID                   string   `yaml:"clientID" json:"clientId"`
	Scopes                     []string `yaml:"scopes" json:"scopes"`
	server.AuthorizationConfig `yaml:",inline"`
	Mediator                   MediatorConfig `yaml:"mediator" json:"mediator"`
}

func (c *AuthorizationConfig) Validate() error {
	if err := c.Mediator.Validate(); err != nil {
		return fmt.Errorf("mediator.%w", err)
	}
	switch c.Mediator.Mode {
	case MediatorMode_UserAgent:
		if c.ClientID == "" {
			return fmt.Errorf("clientID('%v')", c.ClientID)
		}
		if c.Authority == "" {
			return fmt.Errorf("authority('%v')", c.Authority)
		}
		if c.OwnerClaim == "" {
			return fmt.Errorf("ownerClaim('%v')", c.OwnerClaim)
		}
	case MediatorMode_None:
	}
	return nil
}

type Config struct {
	Addr              string                         `yaml:"address" json:"address"`
	EnforcementPolicy server.EnforcementPolicyConfig `yaml:"enforcementPolicy" json:"enforcementPolicy"`
	KeepAlive         server.KeepAliveConfig         `yaml:"keepAlive" json:"keepAlive"`
	TLS               TLSConfig                      `yaml:"tls" json:"tls"`
	Authorization     AuthorizationConfig            `yaml:"authorization" json:"authorization"`
}

func (c *Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("address('%v')", c.Addr)
	}
	if err := c.TLS.Validate(); err != nil {
		return fmt.Errorf("tls.%w", err)
	}

	return nil
}
