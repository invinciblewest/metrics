package senders

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendMetrics(t *testing.T) {
	sender := createSender("http://localhost:8080")

	t.Run("without metrics", func(t *testing.T) {
		st := storage.NewMemStorage()
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

		st := storage.NewMemStorage()
		st.UpdateGauge("testG", 3.14)
		st.UpdateCounter("testC", int64(314))

		err := SendMetrics(st, sender)
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		st := storage.NewMemStorage()
		st.UpdateCounter("testC", int64(314))

		err := SendMetrics(st, sender)
		assert.Error(t, err)

		st.UpdateGauge("testG", 3.14)
		err = SendMetrics(st, sender)
		assert.Error(t, err)
	})

}
