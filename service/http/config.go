package http

import (
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
}

func (c *Config) Validate() error {
	return c.Config.Validate()
}
