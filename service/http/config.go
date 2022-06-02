package http

import (
	"github.com/plgd-dev/hub/v2/pkg/net/listener"
)

type Config struct {
	Connection listener.Config `yaml:",inline" json:",inline"`
}

func (c *Config) Validate() error {
	return c.Connection.Validate()
}
