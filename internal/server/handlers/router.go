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
	r.Use(gzipMiddleware())

	r.Route("/update", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateFromJSON)
		r.Post("/{type}/{name}/{value}", metricsHandler.UpdateFromQuery)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", metricsHandler.GetJSON)
		r.Get("/{type}/{name}", metricsHandler.GetString)
	})

	return r
}
