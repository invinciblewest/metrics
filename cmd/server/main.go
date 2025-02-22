package main

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	st := storage.NewMemStorage()

	err := run(":8080", st)
	if err != nil {
		panic(err)
	}
}

func run(addr string, st storage.Storage) error {
	mux := http.NewServeMux()

	mux.Handle("/update/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateHandler(w, r, st)
	}))

	return http.ListenAndServe(addr, mux)
}

func updateHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	explodedUrl := strings.Split(strings.Trim(r.URL.Path, "/update/"), "/")
	if len(explodedUrl) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := explodedUrl[0], explodedUrl[1], explodedUrl[2]

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
