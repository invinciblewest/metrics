package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	service services.MetricsService
}

func NewHandler(service services.MetricsService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metric := models.Metric{
		ID:    metricName,
		MType: metricType,
	}
	switch metric.MType {
	case models.TypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Value = &value
	case models.TypeCounter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Delta = &delta
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := h.service.Update(ctx, metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metrics.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	updatedMetrics, err := h.service.Update(ctx, metrics)
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

func (h *Handler) UpdateMetricsBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics []models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		logger.Log.Error("failed to decode metrics", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metrics) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBatch(ctx, metrics); err != nil {
		logger.Log.Error("failed to update batch", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	metrics, err := h.service.Get(ctx, metricType, metricName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) || errors.Is(err, storage.ErrWrongType) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metrics.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := h.service.Get(ctx, metrics.MType, metrics.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) ||
			errors.Is(err, storage.ErrWrongType) ||
			errors.Is(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			logger.Log.Error("failed to get metric", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
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

func (h *Handler) PingStorage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.service.PingStorage(ctx) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
