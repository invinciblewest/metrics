package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/server/services"
	"net/http"
	"strconv"
)

type MetricsHandler struct {
	service services.MetricsService
}

func NewMetricsHandler(service services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		service: service,
	}
}

func (h *MetricsHandler) UpdateFromQuery(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metrics := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}
	if !metrics.CheckType() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch metrics.MType {
	case models.TypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics.Value = &value
	case models.TypeCounter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics.Delta = &delta
	}

	if _, err := h.service.Update(metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricsHandler) UpdateFromJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metrics.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	updatedMetrics, err := h.service.Update(metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(updatedMetrics); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *MetricsHandler) GetString(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	metrics, err := h.service.Get(metricType, metricName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metrics.MType {
	case models.TypeGauge:
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(strconv.FormatFloat(*metrics.Value, 'f', -1, 64)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case models.TypeCounter:
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(strconv.FormatInt(*metrics.Delta, 10)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *MetricsHandler) GetJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metrics.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := h.service.Get(metrics.MType, metrics.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
