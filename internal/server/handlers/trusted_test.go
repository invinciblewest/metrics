package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trustedSubnetMiddleware(t *testing.T) {
	const subnet = "192.168.1.0/24"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name          string
		trustedSubnet string
		xRealIP       string
		wantCode      int
	}{
		{
			name:          "no X-Real-IP header",
			trustedSubnet: subnet,
			xRealIP:       "",
			wantCode:      http.StatusForbidden,
		},
		{
			name:          "invalid X-Real-IP",
			trustedSubnet: subnet,
			xRealIP:       "not-an-ip",
			wantCode:      http.StatusForbidden,
		},
		{
			name:          "IP not in subnet",
			trustedSubnet: subnet,
			xRealIP:       "10.0.0.1",
			wantCode:      http.StatusForbidden,
		},
		{
			name:          "IP in subnet",
			trustedSubnet: subnet,
			xRealIP:       "192.168.1.42",
			wantCode:      http.StatusOK,
		},
		{
			name:          "empty trustedSubnet (no check)",
			trustedSubnet: "",
			xRealIP:       "",
			wantCode:      http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := trustedSubnetMiddleware(tt.trustedSubnet)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			rr := httptest.NewRecorder()
			mw(handler).ServeHTTP(rr, req)
			assert.Equal(t, tt.wantCode, rr.Code)
		})
	}
}
