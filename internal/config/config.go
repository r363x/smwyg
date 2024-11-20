package config

import (
	"gopkg.in/yaml.v2"
	"os"
    "fmt"
)

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database"`
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

    fmt.Printf("%#v\n", cfg)

	return &cfg, nil
}
