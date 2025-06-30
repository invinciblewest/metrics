package config

import (
	"flag"

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
}

// GetConfig считывает конфигурацию агента из флагов командной строки и переменных окружения.
func GetConfig() (Config, error) {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&config.PollInterval, "p", 2, "poll interval (sec)")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval (sec)")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.HashKey, "k", "", "hash key")
	flag.IntVar(&config.RateLimit, "L", 2, "rate limit")
	flag.BoolVar(&config.Pprof, "pprof", false, "enable pprof")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
