package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string   `yaml:"lb_port"`
	Backends []string `yaml:"backends"`
	Strategy string   `yaml:"strategy"`
}

func GetLbConfig(path string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Backends) == 0 {
		return nil, errors.New("backend hosts expected, none provided")
	}

	if cfg.Port == "" {
		return nil, errors.New("load balancer port not found")
	}

	if cfg.Strategy == "" {
		cfg.Strategy = "round-robin"
	}

	return &cfg, nil
}
