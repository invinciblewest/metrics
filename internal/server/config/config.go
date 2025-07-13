package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурацию сервера, включая адрес сервера, интервалы опроса и отчета.
type Config struct {
	Address         string `env:"ADDRESS"`           // Адрес сервера, на котором будет запущен сервер.
	LogLevel        string `env:"LOG_LEVEL"`         // Уровень логирования, например, "info", "debug", "error".
	StoreInterval   int    `env:"STORE_INTERVAL"`    // Интервал сохранения метрик в хранилище в секундах.
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // Путь к файлу, в котором будет храниться информация о метриках.
	Restore         bool   `env:"RESTORE"`           // Флаг, указывающий, нужно ли восстанавливать метрики из файла при запуске сервера.
	DatabaseDSN     string `env:"DATABASE_DSN"`      // DSN (Data Source Name) для подключения к базе данных, если используется.
	HashKey         string `env:"KEY"`               // Ключ для хеширования метрик и проверки их целостности.
	CryptoKey       string `env:"CRYPTO_KEY"`        // Приватный ключ для проверки метрик от агента.
}

// GetConfig считывает конфигурацию сервера из флагов командной строки и переменных окружения.
func GetConfig() (Config, error) {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.IntVar(&config.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&config.FileStoragePath, "f", "./storage.json", "storage path")
	flag.BoolVar(&config.Restore, "r", true, "restore")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&config.HashKey, "k", "", "hash key")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "path to crypto key for metrics encryption")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
