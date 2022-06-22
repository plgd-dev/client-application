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

package grpc

import (
	"fmt"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/net/grpc/server"
)

type Config struct {
	server.Config            `yaml:",inline"`
	DefaultGetDevicesRequest pb.GetDevicesRequestConfig `yaml:"defaultGetDevicesRequest" json:"defaultGetDevicesRequest"`
}

func (c *Config) Validate() error {
	if err := c.DefaultGetDevicesRequest.Validate(); err != nil {
		return fmt.Errorf("defaultGetDevicesRequest.%w", err)
	}
	return c.Config.Validate()
}
