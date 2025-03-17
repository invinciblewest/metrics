package main

import (
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
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
	st := storage.NewMemStorage()

	if err := run(cfg.Address, st); err != nil {
		log.Fatal(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)

	logger.Log.Info("Server is starting",
		zap.String("address", addr),
	)
	return http.ListenAndServe(addr, r)
}
