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

func GetMetricHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")

	if mName == "" || mType != "gauge" && mType != "counter" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		v, err := st.GetGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
	default:
		v, err := st.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(strconv.FormatInt(v, 10)))
	}
}
