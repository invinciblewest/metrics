package handlers

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRouter(t *testing.T) {
	st := storage.NewMemStorage()
	r := GetRouter(st)
	assert.Implements(t, (*http.Handler)(nil), r)
}
