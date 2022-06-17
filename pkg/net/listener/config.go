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
