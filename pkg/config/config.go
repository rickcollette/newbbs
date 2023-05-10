package config

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	DBPath string
}

func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		DBPath: cfg.Section("").Key("db_path").String(),
	}, nil
}
