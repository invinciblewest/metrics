package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
	"strconv"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType != "gauge" && metricType != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		st.UpdateGauge(metricName, value)
		return
	default:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		st.UpdateCounter(metricName, value)
		return
	}
}
