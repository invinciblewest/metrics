package handlers

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	explodedURL := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
	if len(explodedURL) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := explodedURL[0], explodedURL[1], explodedURL[2]

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		st.UpdateGauge(metricName, value)
		return
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		st.UpdateCounter(metricName, value)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
