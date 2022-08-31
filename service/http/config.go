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

type Config struct {
	listener.Config `yaml:",inline"`
	Server          server.Config `yaml:",inline" json:",inline"`
	CORS            CORSConfig    `yaml:"cors"  json:"cors"`
	UI              UIConfig      `yaml:"ui" json:"ui"`
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
