package handlers

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		code   int
		target string
	}{
		{
			name:   "method not allowed",
			method: http.MethodGet,
			code:   http.StatusMethodNotAllowed,
			target: "/update/",
		},
		{
			name:   "not found",
			method: http.MethodPost,
			code:   http.StatusNotFound,
			target: "/update/123",
		},
		{
			name:   "gauge error",
			method: http.MethodPost,
			code:   http.StatusBadRequest,
			target: "/update/gauge/test/unknown",
		},
		{
			name:   "gauge success",
			method: http.MethodPost,
			code:   http.StatusOK,
			target: "/update/gauge/test/3.14",
		},
		{
			name:   "counter error",
			method: http.MethodPost,
			code:   http.StatusBadRequest,
			target: "/update/counter/test/unknown",
		},
		{
			name:   "counter success",
			method: http.MethodPost,
			code:   http.StatusOK,
			target: "/update/counter/test/314",
		},
		{
			name:   "unknown type",
			method: http.MethodPost,
			code:   http.StatusBadRequest,
			target: "/update/unknown/test/1",
		},
	}

	st := storage.NewMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.target, nil)
			w := httptest.NewRecorder()

			UpdateMetricHandler(w, r, st)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
