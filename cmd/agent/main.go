package main

import (
	"github.com/invinciblewest/metrics/internal/agent"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/config"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	st := storage.NewMemStorage()
	collectorsList := []collectors.Collector{
		collectors.NewRuntimeCollector(),
	}

	addr := "http://" + cfg.Address
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, http.DefaultClient),
	}

	agent := agent.NewAgent(st, collectorsList, sendersList, cfg.PollInterval, cfg.ReportInterval)
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
