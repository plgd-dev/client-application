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

	"github.com/plgd-dev/client-application/service/config"
)

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
	cfg := config.DefaultConfig(configDirectoryPath + "/www")
	if err := os.WriteFile(configPath, []byte(cfg.String()), 0o600); err != nil {
		return fmt.Errorf("cannot write default config: %w", err)
	}
	os.Args = append(os.Args, "--config", configPath)
	return nil
}
