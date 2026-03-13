package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configer interface {
}

type Config struct {
	Backuper BackuperConfig `yaml:"backuper"`
	path     string
}

type BackuperConfig struct {
	Name          string                    `yaml:"name"`
	Version       float64                   `yaml:"version"`
	Logging       LoggingConfig             `yaml:"logging"`
	Postgres      map[string]PostgresConfig `yaml:"postgres"`
	Mysql         map[string]MysqlConfig    `yaml:"mysql"`
	StorageConfig StorageConfig             `yaml:"storage_config"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.path = path

	return &cfg, nil
}
