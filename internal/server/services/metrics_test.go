package services

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/models"
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
		metricID    string
		metricType  string
		delta       int64
		value       float64
		expectError bool
	}{
		{
			name:        "type error",
			metricID:    "test",
			metricType:  "test",
			expectError: true,
		},
		{
			name:        "gauge success",
			metricID:    "test",
			metricType:  models.TypeGauge,
			value:       3.14,
			expectError: false,
		},
		{
			name:        "counter success",
			metricID:    "test",
			metricType:  models.TypeCounter,
			delta:       314,
			expectError: false,
		},
	}

	st := storage.NewMemStorage()
	service := NewMetricsService(st)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics := models.Metrics{
				ID:    test.metricID,
				MType: test.metricType,
				Delta: &test.delta,
				Value: &test.value,
			}
			fmt.Printf("obj: %+v\n", metrics)
			_, err := service.Update(metrics)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricsService_Get(t *testing.T) {
	expectedValue := 3.14
	expectedDelta := int64(314)

	tests := []struct {
		name          string
		metricID      string
		metricType    string
		expectedDelta int64
		expectedValue float64
		expectError   bool
	}{
		{
			name:        "type error",
			metricID:    "unknown",
			metricType:  "unknown",
			expectError: true,
		},
		{
			name:        "gauge not found",
			metricID:    "unknown",
			metricType:  models.TypeGauge,
			expectError: true,
		},
		{
			name:          "gauge success",
			metricID:      "testG",
			metricType:    models.TypeGauge,
			expectedValue: expectedValue,
			expectError:   false,
		},
		{
			name:        "counter not found",
			metricID:    "unknown",
			metricType:  models.TypeCounter,
			expectError: true,
		},
		{
			name:          "counter success",
			metricID:      "testC",
			metricType:    models.TypeCounter,
			expectedDelta: expectedDelta,
			expectError:   false,
		},
	}

	st := storage.NewMemStorage()
	st.UpdateGauge(models.Metrics{
		ID:    "testG",
		MType: models.TypeGauge,
		Value: &expectedValue,
	})
	st.UpdateCounter(models.Metrics{
		ID:    "testC",
		MType: models.TypeCounter,
		Delta: &expectedDelta,
	})
	service := NewMetricsService(st)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := service.Get(test.metricType, test.metricID)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				switch test.metricType {
				case models.TypeGauge:
					assert.Equal(t, test.expectedValue, *result.Value)
				case models.TypeCounter:
					assert.Equal(t, test.expectedDelta, *result.Delta)
				}
			}
		})
	}
}
