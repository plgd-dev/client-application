package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	service "github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

var Version = "unknown version"

func main() {
	var opts struct {
		Version    bool   `short:"v" long:"version" description:"version"`
		ConfigPath string `long:"config" description:"yaml config file path"`
	}
	_, _ = flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown).Parse()
	if opts.Version {
		fmt.Println(Version)
		return
	}
	if err := resolveDefaultConfig(opts.ConfigPath); err != nil {
		log.Fatalf("cannot create default config: %v", err)
	}
	var cfg service.Config
	if err := config.LoadAndValidateConfig(&cfg); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	if _, err := os.Stat(cfg.APIs.HTTP.UI.Directory); cfg.APIs.HTTP.UI.Enabled && err != nil {
		if err := extractUI(cfg.APIs.HTTP.UI.Directory); err != nil {
			log.Fatalf("cannot extract UI: %v", err)
		}
	}
	logger := log.NewLogger(cfg.Log)
	log.Set(logger)
	log.Debugf("config:\n%v", cfg.String())
	s, err := service.New(context.Background(), cfg, logger)
	if err != nil {
		log.Fatalf("cannot create service: %v", err)
	}
	err = s.Serve()
	if err != nil {
		log.Fatalf("cannot serve service: %v", err)
	}
}
