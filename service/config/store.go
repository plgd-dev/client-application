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

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func SaveToFile(cfg interface{}, f *os.File) error {
	return yaml.NewEncoder(f).Encode(cfg)
}

func Save(cfg interface{}, path string) (err error) {
	var f *os.File
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() {
		if e := f.Close(); e != nil && err == nil {
			err = e
		}
	}()
	err = SaveToFile(cfg, f)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func Store(cfg interface{}, path string) error {
	tmpPath := path + ".tmp"
	if err := Save(cfg, tmpPath); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
