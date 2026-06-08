package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey              string `yaml:"api_key"`
	Model               string `yaml:"model"`
	BaseURL             string `yaml:"base_url"`
	MaxCompletionTokens int    `yaml:"max_completion_tokens"` // LLM 回复的最大 token 数，默认 4096
	MaxRounds           int    `yaml:"max_rounds"`            // Agent 最大循环轮次
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://token-plan-cn.xiaomimimo.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "mimo-v2.5-pro"
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("api_key is required")
	}
	if cfg.MaxCompletionTokens == 0 {
		cfg.MaxCompletionTokens = 4096
	}
	if cfg.MaxRounds == 0 {
		cfg.MaxRounds = 10
	}

	return cfg, nil
}
