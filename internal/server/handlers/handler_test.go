package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/server/services"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler_UpdateFromQuery(t *testing.T) {
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

	st := memstorage.NewMemStorage("", false)
	server := httptest.NewServer(newRouter(st))
	defer server.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = server.URL + test.target

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.code, resp.StatusCode())
		})
	}
}

func TestMetricsHandler_UpdateFromJSON(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		method       string
		target       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method not allowed",
			contentType:  "application/json",
			method:       http.MethodGet,
			target:       "/update/",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "wrong content type",
			contentType:  "application/xml",
			method:       http.MethodPost,
			target:       "/update/",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "invalid body",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/update/",
			body:         `{"wrong": true`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "wrong entity",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/update/",
			body:         `{"wrong": true}`,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "wrong type",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/update/",
			body:         `{"id": "test1", "type": "unknown"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "success",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/update/",
			body:         `{"id":"test","type":"counter","delta":1}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"test","type":"counter","delta":1}`,
		},
	}

	st := memstorage.NewMemStorage("", false)
	server := httptest.NewServer(newRouter(st))
	defer server.Close()
	client := resty.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := client.R()
			req.Method = test.method
			req.URL = server.URL + test.target
			req.SetHeader("Content-Type", test.contentType)
			if test.body != "" {
				req.SetBody(test.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode())
			if test.expectedBody != "" {
				assert.JSONEq(t, test.expectedBody, string(resp.Body()))
			}
		})
	}
}

func TestMetricsHandler_GetString(t *testing.T) {
	ctx := context.TODO()
	st := memstorage.NewMemStorage("", false)
	testG := 3.14
	testC := int64(314)
	err := st.UpdateGauge(ctx, models.Metric{
		ID:    "testG",
		MType: models.TypeGauge,
		Value: &testG,
	})
	assert.NoError(t, err)
	err = st.UpdateCounter(ctx, models.Metric{
		ID:    "testC",
		MType: models.TypeCounter,
		Delta: &testC,
	})
	assert.NoError(t, err)

	server := httptest.NewServer(newRouter(st))
	defer server.Close()

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
			req.URL = server.URL + test.target

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode())
			if test.expectedBody != "" {
				assert.Equal(t, test.expectedBody, string(bytes.TrimRight(resp.Body(), "\n")))
			}
		})
	}
}

func TestMetricsHandler_GetJSON(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		method       string
		target       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method not allowed",
			contentType:  "application/json",
			method:       http.MethodGet,
			target:       "/value/",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "wrong content type",
			contentType:  "application/xml",
			method:       http.MethodPost,
			target:       "/value/",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "invalid body",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/value/",
			body:         `{"wrong": true`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "wrong entity",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/value/",
			body:         `{"wrong": true}`,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "not found",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/value/",
			body:         `{"id": "unknown", "type": "gauge"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "success gauge",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/value/",
			body:         `{"id": "testG", "type": "gauge"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": "testG", "type": "gauge", "value": 3.14}`,
		},
		{
			name:         "success counter",
			contentType:  "application/json",
			method:       http.MethodPost,
			target:       "/value/",
			body:         `{"id": "testC", "type": "counter"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": "testC", "type": "counter", "delta": 314}`,
		},
	}

	ctx := context.TODO()
	st := memstorage.NewMemStorage("", false)
	testG := 3.14
	testC := int64(314)
	err := st.UpdateGauge(ctx, models.Metric{
		ID:    "testG",
		MType: models.TypeGauge,
		Value: &testG,
	})
	assert.NoError(t, err)
	err = st.UpdateCounter(ctx, models.Metric{
		ID:    "testC",
		MType: models.TypeCounter,
		Delta: &testC,
	})
	assert.NoError(t, err)

	server := httptest.NewServer(newRouter(st))
	defer server.Close()
	client := resty.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := client.R()
			req.Method = test.method
			req.URL = server.URL + test.target
			req.SetHeader("Content-Type", test.contentType)
			if test.body != "" {
				req.SetBody(test.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode())
			if test.expectedBody != "" {
				assert.JSONEq(t, test.expectedBody, string(resp.Body()))
			}
		})
	}

}

func newRouter(st storage.Storage) http.Handler {
	return GetRouter(
		NewHandler(
			services.NewMetricsService(st),
		),
		"",
	)
}
