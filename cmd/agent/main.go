package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
	"time"
)

type AgentConfig struct {
	address        string `env:"ADDRESS"`
	pollInterval   int    `env:"POLL_INTERVAL"`
	reportInterval int    `env:"REPORT_INTERVAL"`
}

func main() {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	pollInterval := flag.Int("p", 2, "poll interval (sec)")
	reportInterval := flag.Int("r", 10, "report interval (sec)")
	flag.Parse()

	var config AgentConfig
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	if config.address != "" {
		serverAddr = &config.address
	}
	if config.pollInterval != 0 {
		pollInterval = &config.pollInterval
	}
	if config.reportInterval != 0 {
		reportInterval = &config.reportInterval
	}

	st := storage.NewMemStorage()
	collectorsList := []collectors.Collector{
		collectors.NewRuntimeCollector(),
	}

	addr := "http://" + *serverAddr
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, http.DefaultClient),
	}

	runAgent(st, collectorsList, sendersList, *pollInterval, *reportInterval)
}

func runAgent(
	st *storage.MemStorage,
	c []collectors.Collector,
	s []senders.Sender,
	pInterval int,
	rInterval int,
) {
	for {
		if err := collectors.CollectMetrics(st, c...); err != nil {
			panic(err)
		}

		pc, _ := st.GetCounter("PollCount")
		if ((int(pc) * pInterval) % rInterval) == 0 {
			if err := senders.SendMetrics(st, s...); err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Duration(pInterval) * time.Second)
	}
}
