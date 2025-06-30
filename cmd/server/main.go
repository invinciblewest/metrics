package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/invinciblewest/metrics/internal/storage/pgstorage"
	_ "github.com/lib/pq"
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

	ctx := context.Background()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Log.Fatal("failed to get config", zap.Error(err))
	}
	err = logger.Initialize(cfg.LogLevel)
	if err != nil {
		logger.Log.Fatal("failed to initialize logger", zap.Error(err))
	}

	var st storage.Storage

	if cfg.DatabaseDSN != "" {
		logger.Log.Info("using PostgreSQL storage")

		var db *sql.DB
		db, err = sql.Open("postgres", cfg.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("failed to connect to database", zap.Error(err))
		}

		if err = pgstorage.InstallSchema(db); err != nil {
			logger.Log.Fatal("failed to install schema", zap.Error(err))
		}

		st = pgstorage.NewPGStorage(db)
		defer func(st storage.Storage, ctx context.Context) {
			err = st.Close(ctx)
			if err != nil {
				logger.Log.Fatal("failed to close storage", zap.Error(err))
			}
		}(st, ctx)
	} else {
		logger.Log.Info("using in-memory storage")
		syncSave := cfg.StoreInterval == 0
		st = memstorage.NewMemStorage(cfg.FileStoragePath, syncSave)

		if cfg.Restore {
			if err = st.Load(ctx); err != nil {
				logger.Log.Error("failed to load metrics from file", zap.Error(err))
			}
		}
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ctx.Done():
					if err = st.Save(ctx); err != nil {
						log.Fatal(err)
					}
					return
				case <-ticker.C:
					if err = st.Save(ctx); err != nil {
						log.Fatal(err)
					}
				}
			}
		}()
	}

	handler := handlers.NewHandler(services.NewMetricsService(st))
	router := handlers.GetRouter(handler, cfg.HashKey)

	if err = run(ctx, cfg.Address, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func checkValue(val string) string {
	if val == "" {
		return "N/A"
	}
	return val
}
