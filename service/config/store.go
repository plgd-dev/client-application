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
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
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
