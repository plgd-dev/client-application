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
	"github.com/plgd-dev/hub/v2/pkg/config"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
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
		log.Errorf("cannot create default config: %v", err)
		return
	}
	var cfg service.Config
	if err := config.LoadAndValidateConfig(&cfg); err != nil {
		log.Errorf("cannot load config: %v", err)
		return
	}
	if _, err := os.Stat(cfg.APIs.HTTP.UI.Directory); cfg.APIs.HTTP.UI.Enabled && err != nil {
		if err := extractUI(cfg.APIs.HTTP.UI.Directory); err != nil {
			log.Errorf("cannot extract UI: %v", err)
		}
	}
	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("cannot create file fileWatcher: %v", err)
		return
	}
	defer func() {
		_ = fileWatcher.Close()
	}()
	logger := log.NewLogger(cfg.Log)
	log.Set(logger)
	log.Debugf("config:\n%v", cfg.String())
	s, err := service.New(context.Background(), cfg, fileWatcher, logger)
	if err != nil {
		log.Errorf("cannot create service: %v", err)
		return
	}
	err = s.Serve()
	if err != nil {
		log.Errorf("cannot serve service: %v", err)
		return
	}
}
