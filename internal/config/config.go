package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Host    string        `envconfig:"SERVER_HOST" default:"localhost"`
		Port    string        `envconfig:"PORT" default:"3000"`
		Timeout time.Duration `envconfig:"SERVER_TIMEOUT" default:"15s"`
	}

	StorageCluster struct {
		NodeAPIPort       string `envconfig:"NODE_API_PORT" default:"9000"`
		NodeIdentifier    string `envconfig:"NODE_IDENTIFIER" default:"amazin-object-storage"`
		NetworkIdentifier string `envconfig:"NODE_NETWORK_IDENTIFIER" default:"object-storage_amazin-object-storage"`
	}

	Discovery struct {
		RefreshDuration time.Duration `envconfig:"DISCOVERY_REFRESH" default:"1s"`
	}
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("object-storage", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
