package listener

import (
	"fmt"

	"github.com/plgd-dev/hub/v2/pkg/security/certManager/server"
)

type TLSConfig struct {
	Enabled       bool `yaml:"enabled" json:"enabled"`
	server.Config `yaml:",inline"`
}

func (c *TLSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	return c.Config.Validate()
}

type Config struct {
	Addr string    `yaml:"address" json:"address"`
	TLS  TLSConfig `yaml:"tls" json:"tls"`
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
