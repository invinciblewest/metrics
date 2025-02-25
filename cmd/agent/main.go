package main

import (
	"github.com/invinciblewest/metrics/internal/agent"
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
	cfg := agent.GetConfig()

	st := storage.NewMemStorage()
	collectorsList := []collectors.Collector{
		collectors.NewRuntimeCollector(),
	}

	addr := "http://" + cfg.Address
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, http.DefaultClient),
	}

	runAgent(st, collectorsList, sendersList, cfg.PollInterval, cfg.ReportInterval)
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
