package services

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetricsService(t *testing.T) {
	st := storage.NewMemStorage()
	service := NewMetricsService(st)
	assert.NotNil(t, service)
}

func TestMetricsService_Update(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		expectError bool
	}{
		{
			name:        "type error",
			metricType:  "test",
			metricName:  "test",
			metricValue: "1",
			expectError: true,
		},
		{
			name:        "parse gauge error",
			metricType:  "gauge",
			metricName:  "testG",
			metricValue: "asd",
			expectError: true,
		},
		{
			name:        "gauge success",
			metricType:  "gauge",
			metricName:  "testG",
			metricValue: "3.14",
			expectError: false,
		},
		{
			name:        "parse counter error",
			metricType:  "counter",
			metricName:  "testC",
			metricValue: "asd",
			expectError: true,
		},
		{
			name:        "counter success",
			metricType:  "counter",
			metricName:  "testC",
			metricValue: "314",
			expectError: false,
		},
	}

	st := storage.NewMemStorage()
	service := NewMetricsService(st)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.Update(test.metricType, test.metricName, test.metricValue)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricsService_GetString(t *testing.T) {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "type error",
			metricType:     "unknown",
			metricName:     "unknown",
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "gauge not found",
			metricType:     "gauge",
			metricName:     "unknown",
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "success gauge",
			metricType:     "gauge",
			metricName:     "testG",
			expectedResult: "3.14",
			expectError:    false,
		},
		{
			name:           "counter not found",
			metricType:     "counter",
			metricName:     "unknown",
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "success counter",
			metricType:     "counter",
			metricName:     "testC",
			expectedResult: "314",
			expectError:    false,
		},
	}

	st := storage.NewMemStorage()
	st.UpdateGauge("testG", 3.14)
	st.UpdateCounter("testC", 314)
	service := NewMetricsService(st)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := service.GetString(test.metricType, test.metricName)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}
