package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/invinciblewest/metrics/internal/server/services"
	"net/http"
)

type MetricsHandler struct {
	service services.MetricsService
}

func NewMetricsHandler(service services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		service: service,
	}
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := h.service.Update(metricType, metricName, metricValue); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	metricValue, err := h.service.GetString(metricType, metricName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(metricValue))
}
