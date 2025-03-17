package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

func GetRouter(st storage.Storage) http.Handler {
	r := chi.NewRouter()
	metricsService := services.NewMetricsService(st)
	metricsHandler := NewMetricsHandler(metricsService)

	r.Use(logger.Middleware())
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.Update)
	r.Get("/value/{type}/{name}", metricsHandler.Get)

	return r
}
