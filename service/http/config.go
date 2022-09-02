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

package http

import (
	"fmt"

	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/hub/v2/pkg/net/http/server"
)

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowedOrigins"  json:"allowedOrigins"`
	AllowedHeaders   []string `yaml:"allowedHeaders"  json:"allowedHeaders"`
	AllowedMethods   []string `yaml:"allowedMethods"  json:"allowedMethods"`
	AllowCredentials bool     `yaml:"allowCredentials"  json:"allowCredentials"`
}

type BasicOAuthClient struct {
	Authority string   `yaml:"authority" json:"authority"`
	ClientID  string   `yaml:"clientID" json:"clientId"`
	Audience  string   `yaml:"audience" json:"audience"`
	Scopes    []string `yaml:"scopes" json:"scopes"`
}

func (c *BasicOAuthClient) Validate() error {
	if c.Authority == "" {
		return fmt.Errorf("authority('%v')", c.Authority)
	}
	if c.ClientID == "" {
		return fmt.Errorf("clientID('%v')", c.ClientID)
	}
	return nil
}

type MediatorMode string

const (
	MediatorMode_None      MediatorMode = "none"
	MediatorMode_UserAgent MediatorMode = "userAgent"
)

type UserAgentConfig struct {
	CertificateAuthorityAddress string           `yaml:"certificateAuthorityAddress" json:"certificateAuthorityAddress"`
	WebOAuthClient              BasicOAuthClient `yaml:"webOAuthClient" json:"webOauthClient"`
}

func (c *UserAgentConfig) Validate() error {
	if c.CertificateAuthorityAddress == "" {
		return fmt.Errorf("certificateAuthorityAddress('%v')", c.CertificateAuthorityAddress)
	}
	if err := c.WebOAuthClient.Validate(); err != nil {
		return fmt.Errorf("webOAuthClient.%w", err)
	}
	return nil
}

// WebConfigurationConfig represents web configuration for user interface exposed via getOAuthConfiguration handler
type WebConfigurationConfig struct {
	HTTPGatewayAddress string          `yaml:"-" json:"httpGatewayAddress"`
	MediatorMode       MediatorMode    `yaml:"mediatorMode" json:"mediatorMode"`
	UserAgentConfig    UserAgentConfig `yaml:"userAgent" json:"userAgent"`
}

func (c *WebConfigurationConfig) Validate() error {
	switch c.MediatorMode {
	case MediatorMode_None, "":
		c.MediatorMode = MediatorMode_None
	case MediatorMode_UserAgent:
		if err := c.UserAgentConfig.Validate(); err != nil {
			return fmt.Errorf("userAgent.%w", err)
		}
	}
	return nil
}

type Config struct {
	listener.Config  `yaml:",inline"`
	Server           server.Config          `yaml:",inline" json:",inline"`
	CORS             CORSConfig             `yaml:"cors"  json:"cors"`
	UI               UIConfig               `yaml:"ui" json:"ui"`
	WebConfiguration WebConfigurationConfig `json:"webConfiguration" yaml:"webConfiguration"`
}

type UIConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Directory string `yaml:"directory" json:"directory"`
}

func (c *UIConfig) Validate() error {
	if c.Enabled && c.Directory == "" {
		return fmt.Errorf("directory('%v')", c.Directory)
	}
	return nil
}

func (c *Config) Validate() error {
	if err := c.UI.Validate(); err != nil {
		return fmt.Errorf("ui.%w", err)
	}
	if err := c.WebConfiguration.Validate(); err != nil {
		return fmt.Errorf("webConfiguration.%w", err)
	}
	return c.Config.Validate()
}
