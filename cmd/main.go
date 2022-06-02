package main

import (
	"context"
	"fmt"

	"github.com/jessevdk/go-flags"
	service "github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

var Version = "unknown version"

func main() {
	var opts struct {
		Version bool `short:"v" long:"version" description:"version"`
	}
	_, _ = flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown).Parse()
	if opts.Version {
		fmt.Println(Version)
		return
	}

	var cfg service.Config
	err := config.LoadAndValidateConfig(&cfg)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	logger := log.NewLogger(cfg.Log)
	log.Set(logger)
	log.Infof("config: %v", cfg.String())
	s, err := service.New(context.Background(), cfg, logger)
	if err != nil {
		log.Fatalf("cannot create service: %v", err)
	}
	err = s.Serve()
	if err != nil {
		log.Fatalf("cannot serve service: %v", err)
	}
}
