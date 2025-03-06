package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func GetConfig() (Config, error) {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if config.Address == "" {
		config.Address = *serverAddr
	}

	return config, nil
}
