package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func GetConfig() Config {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	pollInterval := flag.Int("p", 2, "poll interval (sec)")
	reportInterval := flag.Int("r", 10, "report interval (sec)")
	flag.Parse()

	var config Config
	_ = env.Parse(&config)

	if config.Address == "" {
		config.Address = *serverAddr
	}
	if config.PollInterval == 0 {
		config.PollInterval = *pollInterval
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = *reportInterval
	}

	return config
}
