package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Host           string `envconfig:"SERVER_HOST" default:"localhost"`
		Port           string `envconfig:"PORT" default:"3000"`
		TimeoutSeconds uint   `envconfig:"SERVER_TIMEOUT" default:"10"`
	}
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("object-storage", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
