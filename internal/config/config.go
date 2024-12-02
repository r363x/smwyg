package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database"`
    LoggingConfig  LoggingConfig `yaml:"logging"`
}

type LoggingConfig struct {
    LogFile  string `yaml:"log_file"`
    LogLevel string `yaml:"log_level"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func Load() (*Config, error) {
	f, err := os.Open("configs/config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
