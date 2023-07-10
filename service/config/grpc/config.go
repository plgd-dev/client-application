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
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/pkg/net/grpc/server"
	grpcServer "github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
)

type Config = server.Config

type ServiceInformation = pb.BuildInfo

var defaultConfig = Config{
	Addr: ":8081",
	TLS: server.TLSConfig{
		Enabled: false,
	},
	EnforcementPolicy: grpcServer.EnforcementPolicyConfig{
		MinTime: 5 * time.Minute,
	},
}

func DefaultConfig() Config {
	return defaultConfig
}
