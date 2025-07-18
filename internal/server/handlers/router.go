package handlers

import (
	"net/http"

	"github.com/invinciblewest/metrics/pkg/encryption"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/invinciblewest/metrics/internal/logger"
	"go.uber.org/zap"
)

// GetRouter создает и настраивает маршрутизатор Chi с заданным обработчиком.
func GetRouter(handler *Handler, hashKey string, cryptor *encryption.Cryptor, trustedSubnet string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.Middleware())
	r.Use(middleware.Recoverer)
	r.Use(hashMiddleware(hashKey))

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<html><body><h1>Metrics</h1></body></html>"))
		if err != nil {
			logger.Log.Error("main page error", zap.Error(err))
		}
	})

	r.Route("/updates", func(r chi.Router) {
		if cryptor != nil {
			r.Use(encryption.DecryptBodyMiddleware(cryptor))
		}
		r.Use(trustedSubnetMiddleware(trustedSubnet))
		r.Use(gzipMiddleware())
		r.Post("/", handler.UpdateMetricsBatch)
	})
	r.Route("/update", func(r chi.Router) {
		if cryptor != nil {
			r.Use(encryption.DecryptBodyMiddleware(cryptor))
		}
		r.Use(trustedSubnetMiddleware(trustedSubnet))
		r.Use(gzipMiddleware())
		r.Post("/", handler.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", handler.UpdateMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Use(gzipMiddleware())

		r.Post("/", handler.GetMetricJSON)
		r.Get("/{type}/{name}", handler.GetMetric)
	})
	r.Route("/ping", func(r chi.Router) {
		r.Use(gzipMiddleware())
		r.Get("/", handler.PingStorage)
	})

	return r
}
