package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TGKey string `yaml:"tg_key"`
	DBURL string `yaml:"db_url"`
}

func New() *Config {
	slog.Info("Сборка конфига")
	data, err := os.ReadFile("internal/config/config.yml")
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
