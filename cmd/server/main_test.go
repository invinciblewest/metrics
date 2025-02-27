package main

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	st := storage.NewMemStorage()
	cfg := config.GetConfig()

	go func() {
		err := run(cfg.Address, st)
		assert.NoError(t, err)
	}()

	sendRequest(t, cfg.Address)
}

func sendRequest(t *testing.T, addr string) bool {
	target := fmt.Sprintf("http://%s/update/gauge/test/3.14", addr)
	resp, err := http.Post(target, "text/plain", nil)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	}()

	return assert.Equal(t, http.StatusOK, resp.StatusCode)
}
