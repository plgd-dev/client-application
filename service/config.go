package service

import (
	"fmt"

	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

//Config represent application configuration
type Config struct {
	Log  log.Config `yaml:"log" json:"log"`
	APIs APIsConfig `yaml:"apis" json:"apis"`
}

func (c *Config) Validate() error {
	if err := c.APIs.Validate(); err != nil {
		return fmt.Errorf("apis.%w", err)
	}
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log.%w", err)
	}
	return nil
}

type APIsConfig struct {
	HTTP http.Config `yaml:"http" json:"http"`
}

func (c *APIsConfig) Validate() error {
	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("http.%w", err)
	}
	return nil
}

//String return string representation of Config
func (c Config) String() string {
	return config.ToString(c)
}
