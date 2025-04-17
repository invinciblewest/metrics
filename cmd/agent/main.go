package main

import (
	"context"
	"errors"
	"github.com/invinciblewest/metrics/internal/agent"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/config"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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
		collectors.NewGopsutilCollector(st),
	}

	addr := "http://" + cfg.Address
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, cfg.HashKey, http.DefaultClient),
	}

	agentApp := agent.NewAgent(st, collectorsList, sendersList, cfg.PollInterval, cfg.ReportInterval)
	if err = agentApp.Run(ctx, cfg.RateLimit); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Log.Error("agent run error", zap.Error(err))
		} else {
			logger.Log.Info("agent stopped")
		}
	}
}
