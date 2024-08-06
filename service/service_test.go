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

package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestServiceServe(t *testing.T) {
	fmt.Printf("%v\n\n", test.MakeConfig(t))

	shutDown := test.SetUp(t)
	// give some time to start-up
	time.Sleep(time.Second)

	defer shutDown()
}

func TestServiceFailSetup(t *testing.T) {
	cfg := test.MakeConfig(t)
	cfg.APIs.GRPC.Addr = cfg.APIs.HTTP.Addr

	ctx := context.Background()
	logger := log.NewLogger(cfg.Log)
	require.NoError(t, cfg.Validate())

	fileWatcher, err := fsnotify.NewWatcher(logger)
	require.NoError(t, err)
	defer func() {
		_ = fileWatcher.Close()
	}()
	_, err = service.New(ctx, cfg, test.NewServiceInformation().GetBuildInfo(), fileWatcher, logger)
	require.Error(t, err)
}
