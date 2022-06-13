package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pkg/net/grpc/server"
	"github.com/plgd-dev/client-application/pkg/net/listener"
	service "github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/client-application/service/device"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/log"
	grpcServer "github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
)

func createDefaultConfig(uiDirectory string) service.Config {
	return service.Config{
		APIs: service.APIsConfig{
			HTTP: service.HTTPConfig{
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
				},
			},
			GRPC: service.GRPCConfig{
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
		Log: log.MakeDefaultConfig(),
		Clients: service.ClientsConfig{
			Device: device.Config{
				COAP: device.CoapConfig{
					MaxMessageSize: 256 * 1024,
					InactivityMonitor: device.InactivityMonitor{
						Timeout: time.Second * 10,
					},
					BlockwiseTransfer: device.BlockwiseTransferConfig{
						Enabled: true,
						SZX:     "1024",
					},
					TLS: device.TLSConfig{
						SubjectUUID:      uuid.NewString(),
						PreSharedKeyUUID: uuid.NewString(),
					},
				},
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
