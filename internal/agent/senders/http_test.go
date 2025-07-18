package senders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPSender(t *testing.T) {
	s := createSender("http://localhost:8080")
	assert.Implements(t, (*Sender)(nil), s)
}

func TestHTTPSender_Send(t *testing.T) {
	ctx := context.TODO()
	t.Run("error create request", func(t *testing.T) {
		s := createSender("https://%123:8080")
		err := s.SendMetric(ctx, createMetrics())
		assert.Error(t, err)
	})
	t.Run("error send", func(t *testing.T) {
		s := createSender("http://localhost:8080")
		err := s.SendMetric(ctx, createMetrics())
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		st := memstorage.NewMemStorage("", false)
		srv := httptest.NewServer(
			handlers.GetRouter(
				handlers.NewHandler(
					services.NewMetricsService(st),
				),
				"",
				nil,
				"",
			),
		)
		defer srv.Close()
		s := createSender(srv.URL)
		err := s.SendMetric(ctx, createMetrics())
		assert.NoError(t, err)
	})
}

func createMetrics() []models.Metric {
	value := 3.14
	return []models.Metric{
		{
			ID:    "test",
			MType: models.TypeGauge,
			Value: &value,
		},
	}
}

func createSender(addr string) *HTTPSender {
	return NewHTTPSender(addr, "", http.DefaultClient, nil)
}
