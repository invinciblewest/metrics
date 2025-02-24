package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

func GetRouter(st storage.Storage) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, req *http.Request) {
		UpdateMetricHandler(w, req, st)
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, req *http.Request) {
		GetMetricHandler(w, req, st)
	})

	return r
}
