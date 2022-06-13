package service

import (
	"fmt"

	"github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/service/grpc"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

// Config represent application configuration
type Config struct {
	Log     log.Config    `yaml:"log" json:"log"`
	APIs    APIsConfig    `yaml:"apis" json:"apis"`
	Clients ClientsConfig `yaml:"clients" json:"clients"`
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
		return fmt.Errorf("http or grpc must be enabled")
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
