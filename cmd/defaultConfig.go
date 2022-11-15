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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pkg/net/grpc/server"
	"github.com/plgd-dev/client-application/pkg/net/listener"
	"github.com/plgd-dev/client-application/service/config"
	"github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/client-application/service/config/http"
	"github.com/plgd-dev/client-application/service/config/remoteProvisioning"
	"github.com/plgd-dev/hub/v2/pkg/log"
	grpcServer "github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
	httpServer "github.com/plgd-dev/hub/v2/pkg/net/http/server"
)

func createDefaultConfig(uiDirectory string) config.Config {
	logCfg := log.MakeDefaultConfig()
	logCfg.Encoding = "console"
	return config.Config{
		APIs: config.APIsConfig{
			HTTP: config.HTTPConfig{
				Enabled: true,
				Config: http.Config{
					Config: listener.Config{
						Addr: ":8080",
						TLS: listener.TLSConfig{
							Enabled: false,
						},
					},
					CORS: http.CORSConfig{
						AllowedOrigins: []string{"*"},
						AllowedHeaders: []string{"Accept", "Accept-Language", "Accept-Encoding", "Content-Type", "Content-Language", "Content-Length", "Origin", "X-CSRF-Token", "Authorization"},
						AllowedMethods: []string{"GET", "PATCH", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
					},
					UI: http.UIConfig{
						Enabled:   true,
						Directory: uiDirectory,
					},
					Server: httpServer.Config{
						ReadTimeout:       time.Second * 8,
						ReadHeaderTimeout: time.Second * 4,
						WriteTimeout:      time.Second * 16,
						IdleTimeout:       time.Second * 30,
					},
				},
			},
			GRPC: config.GRPCConfig{
				Enabled: true,
				Config: server.Config{
					Addr: ":8081",
					TLS: server.TLSConfig{
						Enabled: false,
					},
					EnforcementPolicy: grpcServer.EnforcementPolicyConfig{
						MinTime: 5 * time.Minute,
					},
				},
			},
		},
		Log: logCfg,
		Clients: config.ClientsConfig{
			Device: device.Config{
				COAP: device.CoapConfig{
					MaxMessageSize: 256 * 1024,
					InactivityMonitor: device.InactivityMonitor{
						Timeout: time.Second * 10,
					},
					BlockwiseTransfer: device.BlockwiseTransferConfig{
						Enabled: true,
						SZXStr:  "1024",
					},
					TLS: device.TLSConfig{
						Authentication: device.AuthenticationPreSharedKey,
						PreSharedKey: device.PreSharedKeyConfig{
							SubjectUUIDStr: uuid.NewString(),
							KeyUUIDStr:     uuid.NewString(),
						},
					},
					OwnershipTransfer: device.OwnershipTransferConfig{
						Methods: []device.OwnershipTransferMethod{device.OwnershipTransferJustWorks},
					},
				},
			},
		},
		RemoteProvisioning: remoteProvisioning.Config{
			Mode: remoteProvisioning.Mode_None,
			UserAgentConfig: remoteProvisioning.UserAgentConfig{
				CSRChallengeStateExpiration: time.Minute * 1,
			},
			Authorization: remoteProvisioning.AuthorizationConfig{
				OwnerClaim: "sub",
			},
		},
	}
}

func resolveDefaultConfig(configPath string) error {
	configPathWasSet := true
	if configPath == "" {
		configPathWasSet = false
		ex, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot get executable path: %w", err)
		}
		exPath := filepath.Dir(ex)
		configPath = exPath + "/config.yaml"
	}
	if _, err := os.Stat(configPath); err == nil {
		if !configPathWasSet {
			os.Args = append(os.Args, "--config", configPath)
		}
		return nil
	}
	configDirectoryPath := filepath.Dir(configPath)
	cfg := createDefaultConfig(configDirectoryPath + "/www")
	if err := os.WriteFile(configPath, []byte(cfg.String()), 0o600); err != nil {
		return fmt.Errorf("cannot write default config: %w", err)
	}
	os.Args = append(os.Args, "--config", configPath)
	return nil
}
