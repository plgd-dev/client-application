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

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plgd-dev/client-application/service/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	configDir, err := os.MkdirTemp("", "test.*****")
	require.NoError(t, err)
	defer func() {
		err = os.RemoveAll(configDir)
		require.NoError(t, err)
	}()
	configPath := filepath.Join(configDir, "config.yaml")
	cfg := config.DefaultConfig(filepath.Dir(configPath) + "/www")
	err = cfg.Validate()
	require.NoError(t, err)
	cfg.SetConfigPath(configPath)
	err = cfg.Store()
	require.NoError(t, err)
	os.Args = append(os.Args, "--config", configPath)
	cfg1, err := config.New(configPath)
	require.NoError(t, err)
	require.Equal(t, cfg, cfg1)
	require.NotEmpty(t, cfg1.String())
	cfg.Clients.Device.COAP.TLS.PreSharedKey.Key = ""
	cfg.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr = ""
	err = cfg.Validate()
	require.NoError(t, err)
	err = cfg.Store()
	require.NoError(t, err)
	cfg2, err := config.New(configPath)
	require.NoError(t, err)
	require.Equal(t, cfg, cfg2)
	require.Empty(t, cfg2.Clients.Device.COAP.TLS.PreSharedKey.Key)
	require.Empty(t, cfg2.Clients.Device.COAP.TLS.PreSharedKey.SubjectIDStr)
}
