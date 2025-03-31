package senders

import (
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendMetrics(t *testing.T) {
	sender := createSender("http://localhost:8080")

	t.Run("without metrics", func(t *testing.T) {
		st := storage.NewMemStorage("", false)
		err := SendMetrics(st, sender)
		assert.NoError(t, err)
	})
	t.Run("success", func(t *testing.T) {
		server := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})}
		go func() {
			_ = server.ListenAndServe()
		}()
		defer func() {
			_ = server.Close()
		}()

		st := storage.NewMemStorage("", false)
		err := st.UpdateGauge(createGaugeMetric())
		assert.NoError(t, err)
		err = st.UpdateCounter(createCounterMetric())
		assert.NoError(t, err)

		err = SendMetrics(st, sender)
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		st := storage.NewMemStorage("", false)
		err := st.UpdateCounter(createCounterMetric())
		assert.NoError(t, err)

		err = SendMetrics(st, sender)
		assert.Error(t, err)

		err = st.UpdateGauge(createGaugeMetric())
		assert.NoError(t, err)
		err = SendMetrics(st, sender)
		assert.Error(t, err)
	})
}

func createCounterMetric() models.Metric {
	delta := int64(314)
	return models.Metric{
		ID:    "testC",
		MType: models.TypeCounter,
		Delta: &delta,
	}
}

func createGaugeMetric() models.Metric {
	value := 3.14
	return models.Metric{
		ID:    "testG",
		MType: models.TypeGauge,
		Value: &value,
	}
}
