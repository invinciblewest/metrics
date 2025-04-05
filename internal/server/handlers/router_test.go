package handlers

import (
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRouter(t *testing.T) {
	st := memstorage.NewMemStorage("", false)
	r := newRouter(st)
	assert.Implements(t, (*http.Handler)(nil), r)
}
