package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func GetConfig() (Config, error) {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&config.PollInterval, "p", 2, "poll interval (sec)")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval (sec)")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
