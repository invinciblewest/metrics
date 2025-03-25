package main

import (
	"context"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	st := storage.NewMemStorage(cfg.FileStoragePath, cfg.StoreInterval == 0)

	if cfg.Restore {
		if err = st.Load(); err != nil {
			log.Fatal(err)
		}
	}

	if cfg.StoreInterval > 0 {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		go func() {
			for {
				select {
				case <-sig:
					if err = st.Save(); err != nil {
						log.Fatal(err)
					}
					return
				case <-ticker.C:
					if err = st.Save(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}()
	}

	if err := run(cfg.Address, st); err != nil {
		log.Fatal(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Log.Info("server is shutting down...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Fatal("server shutdown error", zap.Error(err))
		}
	}()

	logger.Log.Info("server is starting",
		zap.String("address", addr),
	)
	return server.ListenAndServe()
}
