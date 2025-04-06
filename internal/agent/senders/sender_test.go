package senders

import (
	"context"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendMetrics(t *testing.T) {
	ctx := context.TODO()
	sender := createSender("http://localhost:8080")

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

		st := memstorage.NewMemStorage("", false)
		err := st.UpdateGauge(ctx, createGaugeMetric())
		assert.NoError(t, err)
		err = st.UpdateCounter(ctx, createCounterMetric())
		assert.NoError(t, err)

		err = SendMetrics(ctx, st, sender)
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		st := memstorage.NewMemStorage("", false)
		err := st.UpdateCounter(ctx, createCounterMetric())
		assert.NoError(t, err)

		err = SendMetrics(ctx, st, sender)
		assert.Error(t, err)

		err = st.UpdateGauge(ctx, createGaugeMetric())
		assert.NoError(t, err)
		err = SendMetrics(ctx, st, sender)
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
