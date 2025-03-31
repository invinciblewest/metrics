package main

import (
	"context"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	syncSave := cfg.StoreInterval == 0
	st := storage.NewMemStorage(cfg.FileStoragePath, syncSave)

	if cfg.Restore {
		if err = st.Load(); err != nil {
			log.Fatal(err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ctx.Done():
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

	router := handlers.GetRouter(
		handlers.NewHandler(
			services.NewMetricsService(st),
		),
	)

	if err := run(ctx, cfg.Address, router); err != nil {
		logger.Log.Fatal("server error", zap.Error(err))
	}
}

func run(ctx context.Context, addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Log.Info("server is shutting down...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Fatal("server shutdown error", zap.Error(err))
		}
	}()

	logger.Log.Info("server is starting", zap.String("address", addr))
	return server.ListenAndServe()
}
