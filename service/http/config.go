package http

import (
	"fmt"

	"github.com/plgd-dev/client-application/pkg/net/listener"
)

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowedOrigins"  json:"allowedOrigins"`
	AllowedHeaders   []string `yaml:"allowedHeaders"  json:"allowedHeaders"`
	AllowedMethods   []string `yaml:"allowedMethods"  json:"allowedMethods"`
	AllowCredentials bool     `yaml:"allowCredentials"  json:"allowCredentials"`
}

type Config struct {
	listener.Config `yaml:",inline"`
	CORS            CORSConfig `yaml:"cors"  json:"cors"`
	UI              UIConfig   `yaml:"ui" json:"ui"`
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
	return c.Config.Validate()
}
