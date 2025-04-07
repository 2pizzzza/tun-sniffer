package config

import (
	yaml "gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TunName string `yaml:"tun_name"`
	TunIP   string `yaml:"tun_ip"`
	TunMask string `yaml:"tun_mask"`
	LogPath string `yaml:"log_path"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
