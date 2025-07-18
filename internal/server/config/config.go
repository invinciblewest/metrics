package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурацию сервера
type Config struct {
	Address         string `env:"ADDRESS"`           // Адрес сервера, на котором будет запущен сервер.
	LogLevel        string `env:"LOG_LEVEL"`         // Уровень логирования, например, "info", "debug", "error".
	StoreInterval   int    `env:"STORE_INTERVAL"`    // Интервал сохранения метрик в хранилище в секундах.
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // Путь к файлу, в котором будет храниться информация о метриках.
	Restore         bool   `env:"RESTORE"`           // Флаг, указывающий, нужно ли восстанавливать метрики из файла при запуске сервера.
	DatabaseDSN     string `env:"DATABASE_DSN"`      // DSN (Data Source Name) для подключения к базе данных, если используется.
	HashKey         string `env:"KEY"`               // Ключ для хеширования метрик и проверки их целостности.
	CryptoKey       string `env:"CRYPTO_KEY"`        // Приватный ключ для проверки метрик от агента.
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`    // CIDR доверенной подсети.
}

// JSONConfig представляет структуру JSON файла конфигурации сервера
type JSONConfig struct {
	Address       string `json:"address"`
	Restore       *bool  `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
}

// GetConfig считывает конфигурацию сервера из флагов командной строки и переменных окружения.
func GetConfig() (Config, error) {
	var config Config
	var configFile string

	config = Config{
		Address:         "localhost:8080",
		LogLevel:        "info",
		StoreInterval:   300,
		FileStoragePath: "./storage.json",
		Restore:         true,
		DatabaseDSN:     "",
		HashKey:         "",
		CryptoKey:       "",
		TrustedSubnet:   "",
	}

	flag.StringVar(&config.Address, "a", config.Address, "server address")
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "log level")
	flag.IntVar(&config.StoreInterval, "i", config.StoreInterval, "store interval")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "storage path")
	flag.BoolVar(&config.Restore, "r", config.Restore, "restore")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "database dsn")
	flag.StringVar(&config.HashKey, "k", config.HashKey, "hash key")
	flag.StringVar(&config.CryptoKey, "crypto-key", config.CryptoKey, "path to crypto key for metrics encryption")
	flag.StringVar(&configFile, "c", "", "config file path")
	flag.StringVar(&configFile, "config", "", "config file path")
	flag.StringVar(&config.TrustedSubnet, "t", config.TrustedSubnet, "trusted subnet in CIDR format")
	flag.Parse()

	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	if configFile != "" {
		jsonConfig, err := loadJSONConfig(configFile)
		if err != nil {
			return Config{}, err
		}
		applyJSONConfig(&config, jsonConfig)
	}

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func loadJSONConfig(filename string) (*JSONConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var jsonConfig JSONConfig
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return nil, err
	}

	return &jsonConfig, nil
}

func applyJSONConfig(config *Config, jsonConfig *JSONConfig) {
	if jsonConfig.Address != "" {
		config.Address = jsonConfig.Address
	}
	if jsonConfig.Restore != nil {
		config.Restore = *jsonConfig.Restore
	}
	if jsonConfig.StoreInterval != "" {
		if duration, err := time.ParseDuration(jsonConfig.StoreInterval); err == nil {
			config.StoreInterval = int(duration.Seconds())
		}
	}
	if jsonConfig.StoreFile != "" {
		config.FileStoragePath = jsonConfig.StoreFile
	}
	if jsonConfig.DatabaseDSN != "" {
		config.DatabaseDSN = jsonConfig.DatabaseDSN
	}
	if jsonConfig.CryptoKey != "" {
		config.CryptoKey = jsonConfig.CryptoKey
	}
	if jsonConfig.TrustedSubnet != "" {
		config.TrustedSubnet = jsonConfig.TrustedSubnet
	}
}
