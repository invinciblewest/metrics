package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetrics_CheckType(t *testing.T) {
	tests := []struct {
		name           string
		mType          string
		expectedResult bool
	}{
		{
			name:           "counter true",
			mType:          TypeCounter,
			expectedResult: true,
		},
		{
			name:           "gauge true",
			mType:          TypeGauge,
			expectedResult: true,
		},
		{
			name:           "unknown false",
			mType:          "unknown",
			expectedResult: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics := Metrics{
				MType: test.mType,
			}
			assert.Equal(t, test.expectedResult, metrics.CheckType())
		})
	}
}
