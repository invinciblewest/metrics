package main

import (
	"context"
	"errors"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/invinciblewest/metrics/internal/storage/pgstorage"
	_ "github.com/lib/pq"
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
	memStorage := memstorage.NewMemStorage(cfg.FileStoragePath, syncSave)
	pgStorage := pgstorage.NewPGStorage(cfg.DatabaseDSN)
	defer pgStorage.Close()

	if cfg.Restore {
		if err = memStorage.Load(); err != nil {
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
					if err = memStorage.Save(); err != nil {
						log.Fatal(err)
					}
					return
				case <-ticker.C:
					if err = memStorage.Save(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}()
	}

	handler := handlers.NewHandler(services.NewMetricsService(memStorage))
	router := handlers.GetRouter(handler)
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if err := pgStorage.Ping(); err != nil {
			logger.Log.Error("ping error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	if err := run(ctx, cfg.Address, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
