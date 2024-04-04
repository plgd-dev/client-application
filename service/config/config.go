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

package config

import (
	"errors"
	"fmt"

	"github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/client-application/service/config/grpc"
	"github.com/plgd-dev/client-application/service/config/http"
	"github.com/plgd-dev/client-application/service/config/remoteProvisioning"
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

// Config represent application configuration
type Config struct {
	Log                log.Config                 `yaml:"log" json:"log"`
	APIs               APIsConfig                 `yaml:"apis" json:"apis"`
	Clients            ClientsConfig              `yaml:"clients" json:"clients"`
	RemoteProvisioning *remoteProvisioning.Config `yaml:"remoteProvisioning" json:"remoteProvisioning"`
	configPath         string                     `yaml:"-" json:"-"`
}

func New(configPath string) (Config, error) {
	if configPath == "" {
		return Config{}, errors.New("path to config is empty")
	}
	var cfg Config
	if err := config.LoadAndValidateConfig(&cfg); err != nil {
		return Config{}, fmt.Errorf("cannot load config: %w", err)
	}
	cfg.configPath = configPath
	return cfg, nil
}

func (c *Config) SetConfigPath(configPath string) {
	c.configPath = configPath
}

func (c *Config) Validate() error {
	if err := c.APIs.Validate(); err != nil {
		return fmt.Errorf("apis.%w", err)
	}
	if err := c.Clients.Validate(); err != nil {
		return fmt.Errorf("clients.%w", err)
	}
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log.%w", err)
	}
	if err := c.RemoteProvisioning.Validate(); err != nil {
		return fmt.Errorf("remoteProvisioning.%w", err)
	}
	return nil
}

type HTTPConfig struct {
	Enabled     bool `yaml:"enabled" json:"enabled"`
	http.Config `yaml:",inline" json:",inline"`
}

type GRPCConfig struct {
	Enabled     bool `yaml:"enabled" json:"enabled"`
	grpc.Config `yaml:",inline" json:",inline"`
}

type APIsConfig struct {
	HTTP HTTPConfig `yaml:"http" json:"http"`
	GRPC GRPCConfig `yaml:"grpc" json:"grpc"`
}

func (c *APIsConfig) Validate() error {
	if !c.HTTP.Enabled && !c.GRPC.Enabled {
		return errors.New("http or grpc must be enabled")
	}
	if c.HTTP.Enabled {
		if err := c.HTTP.Validate(); err != nil {
			return fmt.Errorf("http.%w", err)
		}
	}
	if c.GRPC.Enabled {
		if err := c.GRPC.Validate(); err != nil {
			return fmt.Errorf("grpc.%w", err)
		}
	}
	return nil
}

type ClientsConfig struct {
	Device device.Config `yaml:"device" json:"device"`
}

func (c *ClientsConfig) Validate() error {
	if err := c.Device.Validate(); err != nil {
		return fmt.Errorf("device.%w", err)
	}
	return nil
}

// String return string representation of Config
func (c Config) String() string {
	return config.ToString(c)
}

func (c Config) Store() error {
	return Store(c, c.configPath)
}

func DefaultConfig(directory string) Config {
	logCfg := log.MakeDefaultConfig()
	logCfg.Encoding = "console"
	return Config{
		Log: logCfg,
		APIs: APIsConfig{
			HTTP: HTTPConfig{
				Enabled: true,
				Config:  http.DefaultConfig(directory),
			},
			GRPC: GRPCConfig{
				Enabled: true,
				Config:  grpc.DefaultConfig(directory),
			},
		},
		Clients: ClientsConfig{
			Device: device.DefaultConfig(),
		},
		RemoteProvisioning: remoteProvisioning.DefaultConfig(),
	}
}
