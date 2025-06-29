package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	APIKey    string `yaml:"api_key"`
	Receivers []int64 `yaml:"receivers"`
}

type Config struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(home, ".config", "gotify", "config.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{Services: make(map[string]ServiceConfig)}, nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Services == nil {
		cfg.Services = make(map[string]ServiceConfig)
	}

	return &cfg, nil
}

func (c *Config) Save(configPath string) error {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config", "gotify")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
		configPath = filepath.Join(configDir, "config.yaml")
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
