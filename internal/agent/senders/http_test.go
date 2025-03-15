package senders

import (
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
		err := s.Send("gauge", "test", "3.14")
		assert.Error(t, err)
	})
	t.Run("error send", func(t *testing.T) {
		s := createSender("http://localhost:8080")
		err := s.Send("gauge", "test", "3.14")
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		st := storage.NewMemStorage()
		srv := httptest.NewServer(handlers.GetRouter(st))
		defer srv.Close()
		s := createSender(srv.URL)
		err := s.Send("gauge", "test", "3.14")
		assert.NoError(t, err)
	})
}

func createSender(addr string) *HTTPSender {
	return NewHTTPSender(addr, http.DefaultClient)
}
