package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func GetConfig() (Config, error) {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	logLevel := flag.String("l", "info", "log level")
	flag.Parse()

	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if config.Address == "" {
		config.Address = *serverAddr
	}

	if config.LogLevel == "" {
		config.LogLevel = *logLevel
	}

	return config, nil
}
