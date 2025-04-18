package main

import (
	"github.com/invinciblewest/metrics/internal/agent"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/config"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	st := memstorage.NewMemStorage("", false)
	collectorsList := []collectors.Collector{
		collectors.NewRuntimeCollector(st),
	}

	addr := "http://" + cfg.Address
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, http.DefaultClient),
	}

	agentApp := agent.NewAgent(st, collectorsList, sendersList, cfg.PollInterval, cfg.ReportInterval)
	if err := agentApp.Run(); err != nil {
		log.Fatal(err)
	}
}
