package handlers

import (
	"github.com/go-resty/resty/v2"
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
			target: "/update/gauge/test/3.14",
		},
		{
			name:   "not found",
			method: http.MethodPost,
			code:   http.StatusNotFound,
			target: "/update/gauge//3.14",
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
	s := httptest.NewServer(GetRouter(st))
	defer s.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = s.URL + test.target

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.code, resp.StatusCode())
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	st := storage.NewMemStorage()
	st.UpdateGauge("testG", 3.14)
	st.UpdateCounter("testC", 314)

	s := httptest.NewServer(GetRouter(st))
	defer s.Close()

	tests := []struct {
		name         string
		target       string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "type not found",
			target:       "/value/unknown/test",
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "gauge not found",
			target:       "/value/gauge/unknown",
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "counter not found",
			target:       "/value/counter/unknown",
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "gauge success",
			target:       "/value/gauge/testG",
			expectedCode: http.StatusOK,
			expectedBody: "3.14",
		},
		{
			name:         "counter success",
			target:       "/value/counter/testC",
			expectedCode: http.StatusOK,
			expectedBody: "314",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = s.URL + test.target

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode())
			if test.expectedBody != "" {
				assert.Equal(t, test.expectedBody, string(resp.Body()))
			}
		})
	}
}
