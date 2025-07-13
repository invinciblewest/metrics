package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурацию агента, включая адрес сервера, интервалы опроса и отчета.
type Config struct {
	Address        string `env:"ADDRESS"`         // Адрес сервера, на который будет отправлять метрики агент.
	PollInterval   int    `env:"POLL_INTERVAL"`   // Интервал опроса метрик в секундах.
	ReportInterval int    `env:"REPORT_INTERVAL"` // Интервал отправки отчетов на сервер в секундах.
	LogLevel       string `env:"LOG_LEVEL"`       // Уровень логирования, например, "info", "debug", "error".
	HashKey        string `env:"KEY"`             // Ключ для хеширования метрик перед отправкой на сервер.
	RateLimit      int    `env:"RATE_LIMIT"`      // Ограничение скорости отправки метрик на сервер (количество метрик в секунду).
	Pprof          bool   `env:"PPROF"`           // Флаг, указывающий, нужно ли включать pprof для профилирования производительности.
	CryptoKey      string `env:"CRYPTO_KEY"`      // Ключ для шифрования метрик перед отправкой на сервер.
}

// JSONConfig представляет структуру JSON файла конфигурации агента
type JSONConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

// GetConfig считывает конфигурацию агента из флагов командной строки и переменных окружения.
func GetConfig() (Config, error) {
	var config Config
	var configFile string

	config = Config{
		Address:        "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		LogLevel:       "info",
		HashKey:        "",
		RateLimit:      2,
		Pprof:          false,
		CryptoKey:      "",
	}

	flag.StringVar(&config.Address, "a", config.Address, "server address")
	flag.IntVar(&config.PollInterval, "p", config.PollInterval, "poll interval (sec)")
	flag.IntVar(&config.ReportInterval, "r", config.ReportInterval, "report interval (sec)")
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "log level")
	flag.StringVar(&config.HashKey, "k", config.HashKey, "hash key")
	flag.IntVar(&config.RateLimit, "L", config.RateLimit, "rate limit")
	flag.BoolVar(&config.Pprof, "pprof", config.Pprof, "enable pprof")
	flag.StringVar(&config.CryptoKey, "crypto-key", config.CryptoKey, "path to crypto key for metrics encryption")
	flag.StringVar(&configFile, "c", "", "config file path")
	flag.StringVar(&configFile, "config", "", "config file path")
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
	if jsonConfig.ReportInterval != "" {
		if duration, err := time.ParseDuration(jsonConfig.ReportInterval); err == nil {
			config.ReportInterval = int(duration.Seconds())
		}
	}
	if jsonConfig.PollInterval != "" {
		if duration, err := time.ParseDuration(jsonConfig.PollInterval); err == nil {
			config.PollInterval = int(duration.Seconds())
		}
	}
	if jsonConfig.CryptoKey != "" {
		config.CryptoKey = jsonConfig.CryptoKey
	}
}
