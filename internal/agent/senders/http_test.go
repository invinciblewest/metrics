package senders

import (
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPSender(t *testing.T) {
	s := createSender("http://localhost:8080")
	assert.Implements(t, (*Sender)(nil), s)
}

func TestHTTPSender_Send(t *testing.T) {
	t.Run("error create request", func(t *testing.T) {
		s := createSender("https://%123:8080")
		err := s.Send(createMetrics())
		assert.Error(t, err)
	})
	t.Run("error send", func(t *testing.T) {
		s := createSender("http://localhost:8080")
		err := s.Send(createMetrics())
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		st := storage.NewMemStorage("", false)
		srv := httptest.NewServer(handlers.GetRouter(st))
		defer srv.Close()
		s := createSender(srv.URL)
		err := s.Send(createMetrics())
		assert.NoError(t, err)
	})
}

func createMetrics() models.Metrics {
	value := 3.14
	return models.Metrics{
		ID:    "test",
		MType: models.TypeGauge,
		Value: &value,
	}
}

func createSender(addr string) *HTTPSender {
	return NewHTTPSender(addr, http.DefaultClient)
}
