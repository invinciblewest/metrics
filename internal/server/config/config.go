package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func GetConfig() (Config, error) {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.IntVar(&config.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&config.FileStoragePath, "f", "./storage.json", "storage path")
	flag.BoolVar(&config.Restore, "r", true, "restore")
	flag.StringVar(&config.DatabaseDSN, "d", "postgresql://root:secret@localhost:54321/metrics?sslmode=disable", "database dsn")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
