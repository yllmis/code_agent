package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
	BaseURL string `yaml:"base_url"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://token-plan-cn.xiaomimimo.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "mimo-v2.5-pro"
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	return &cfg, nil
}
