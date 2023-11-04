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
	"context"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	service "github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/client-application/service/config"
	"github.com/plgd-dev/client-application/service/config/grpc"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

var (
	Version    = "unknown version"
	BuildDate  = "unknown date"
	CommitHash = "unknown hash"
	CommitDate = "unknown commit date"
	ReleaseURL = "unknown url"
)

func loadConfig() config.Config {
	var opts struct {
		Version    bool   `short:"v" long:"version" description:"version"`
		ConfigPath string `long:"config" description:"yaml config file path"`
	}
	_, _ = flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown).Parse()
	if opts.Version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if err := resolveDefaultConfig(opts.ConfigPath); err != nil {
		log.Errorf("cannot create default config: %v", err)
		os.Exit(1)
	}
	// parse line arguments again because resolveDefaultConfig can set config path
	_, _ = flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown).Parse()
	cfg, err := config.New(opts.ConfigPath)
	if err != nil {
		log.Errorf("cannot load config: %v", err)
		os.Exit(1)
	}
	if _, err = os.Stat(cfg.APIs.HTTP.UI.Directory); cfg.APIs.HTTP.UI.Enabled && err != nil {
		if err = extractUI(cfg.APIs.HTTP.UI.Directory); err != nil {
			log.Errorf("cannot extract UI: %v", err)
			os.Exit(1)
		}
	}
	if cfg.APIs.HTTP.Enabled && cfg.APIs.HTTP.TLS.Enabled &&
		!checkSelfSignedCertificate(cfg.APIs.HTTP.TLS.CertFile.FilePath(), cfg.APIs.HTTP.TLS.KeyFile.FilePath()) {
		if err = generateSelfSigned(cfg.APIs.HTTP.TLS.CertFile.FilePath(), cfg.APIs.HTTP.TLS.KeyFile.FilePath()); err != nil {
			log.Errorf("cannot generate self signed certificate for HTTP: %v", err)
			os.Exit(1)
		}
	}
	if cfg.APIs.GRPC.Enabled && cfg.APIs.GRPC.TLS.Enabled &&
		!checkSelfSignedCertificate(cfg.APIs.GRPC.TLS.CertFile.FilePath(), cfg.APIs.GRPC.TLS.KeyFile.FilePath()) {
		if err = generateSelfSigned(cfg.APIs.GRPC.TLS.CertFile.FilePath(), cfg.APIs.GRPC.TLS.KeyFile.FilePath()); err != nil {
			log.Errorf("cannot generate self signed certificate for GRPC: %v", err)
			os.Exit(1)
		}
	}
	return cfg
}

func main() {
	cfg := loadConfig()
	logger := log.NewLogger(cfg.Log)
	log.Set(logger)
	fileWatcher, err := fsnotify.NewWatcher(logger)
	if err != nil {
		log.Errorf("cannot create file fileWatcher: %v", err)
		os.Exit(1)
	}
	log.Debugf("version: %v, buildDate: %v, buildRevision %v", Version, BuildDate, CommitHash)
	log.Debugf("config:\n%v", cfg.String())
	info := grpc.ServiceInformation{
		Version:    Version,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
		CommitDate: CommitDate,
		ReleaseUrl: ReleaseURL,
	}

	s, err := service.New(context.Background(), cfg, &info, fileWatcher, logger)
	if err != nil {
		log.Errorf("cannot create service: %v", err)
		os.Exit(1)
	}
	err = s.Serve()
	if err != nil {
		log.Errorf("cannot serve service: %v", err)
		os.Exit(1)
	}
}
