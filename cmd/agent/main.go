package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/invinciblewest/metrics/pkg/encryption"

	"github.com/invinciblewest/metrics/internal/agent"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/config"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"go.uber.org/zap"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	fmt.Println("Build version:", checkValue(BuildVersion))
	fmt.Println("Build date:   ", checkValue(BuildDate))
	fmt.Println("Build commit: ", checkValue(BuildCommit))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Pprof {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
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

	var cryptor *encryption.Cryptor
	if cfg.CryptoKey != "" {
		cryptor, err = encryption.NewCryptor(cfg.CryptoKey, "")
		if err != nil {
			logger.Log.Fatal("failed to initialize cryptor", zap.Error(err))
		}
	}

	addr := "http://" + cfg.Address
	sendersList := []senders.Sender{
		senders.NewHTTPSender(addr, cfg.HashKey, http.DefaultClient, cryptor),
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

func checkValue(val string) string {
	if val == "" {
		return "N/A"
	}
	return val
}
