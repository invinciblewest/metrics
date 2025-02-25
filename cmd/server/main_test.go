package main

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	st := storage.NewMemStorage()
	addr := ":8081"

	go func() {
		err := run(addr, st)
		assert.NoError(t, err)
	}()

	sendRequest(t, addr)
}

func TestMainFunction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		go main()
		require.True(t, sendRequest(t, ":8080"))
	})
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			main()
		})
	})
}

func sendRequest(t *testing.T, addr string) bool {
	target := fmt.Sprintf("http://localhost%s/update/gauge/test/3.14", addr)
	resp, err := http.Post(target, "text/plain", nil)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	}()

	return assert.Equal(t, http.StatusOK, resp.StatusCode)
}
