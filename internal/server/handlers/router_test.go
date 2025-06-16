package handlers

import (
	"net/http"
	"testing"

	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestGetRouter(t *testing.T) {
	st := memstorage.NewMemStorage("", false)
	r := newRouter(st)
	assert.Implements(t, (*http.Handler)(nil), r)
}
