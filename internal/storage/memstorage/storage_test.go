package memstorage

import (
	"context"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	st := NewMemStorage("", false)
	assert.Implements(t, (*storage.Storage)(nil), st)
}

func TestMemStorage_Gauge(t *testing.T) {
	ctx := context.TODO()
	st := NewMemStorage("", false)

	f1 := 3.14
	f2 := 14.3
	list := storage.GaugeList{
		"test1": models.Metric{
			ID:    "test1",
			MType: models.TypeGauge,
			Value: &f1,
		},
		"test2": models.Metric{
			ID:    "test2",
			MType: models.TypeGauge,
			Value: &f2,
		},
	}
	t.Run("update gauge", func(t *testing.T) {
		for _, v := range list {
			err := st.UpdateGauge(ctx, v)
			assert.NoError(t, err)
		}
	})
	t.Run("update gauge error", func(t *testing.T) {
		err := st.UpdateGauge(ctx, models.Metric{
			ID:    "test3",
			MType: models.TypeCounter,
		})
		assert.Error(t, err)
	})
	t.Run("get gauge", func(t *testing.T) {
		for k, v := range list {
			val, err := st.GetGauge(ctx, k)
			assert.NoError(t, err)
			assert.Equal(t, v.Value, val.Value)
		}
	})
	t.Run("get gauge error", func(t *testing.T) {
		_, err := st.GetGauge(ctx, "unknown")
		assert.Error(t, err)
	})
	t.Run("get gauge list", func(t *testing.T) {
		assert.Equal(t, list, st.GetGaugeList(ctx))
	})
}

func TestMemStorage_Counter(t *testing.T) {
	ctx := context.TODO()
	st := NewMemStorage("", false)

	c1 := int64(1)
	c2 := int64(2)

	list := storage.CounterList{
		"test1": models.Metric{
			ID:    "test1",
			MType: models.TypeCounter,
			Delta: &c1,
		},
		"test2": models.Metric{
			ID:    "test2",
			MType: models.TypeCounter,
			Delta: &c2,
		},
	}
	t.Run("update counter", func(t *testing.T) {
		for _, v := range list {
			err := st.UpdateCounter(ctx, v)
			assert.NoError(t, err)
		}
	})
	t.Run("update counter error", func(t *testing.T) {
		err := st.UpdateCounter(ctx, models.Metric{
			ID:    "test3",
			MType: models.TypeGauge,
		})
		assert.Error(t, err)
	})
	t.Run("get counter", func(t *testing.T) {
		for k, v := range list {
			val, err := st.GetCounter(ctx, k)
			assert.NoError(t, err)
			assert.Equal(t, v.Delta, val.Delta)
		}
	})
	t.Run("get counter error", func(t *testing.T) {
		_, err := st.GetCounter(ctx, "unknown")
		assert.Error(t, err)
	})
	t.Run("get counter list", func(t *testing.T) {
		assert.Equal(t, list, st.GetCounterList(ctx))
	})
	t.Run("increment counter", func(t *testing.T) {
		for _, v := range list {
			err := st.UpdateCounter(ctx, v)
			assert.NoError(t, err)
		}
	})
	t.Run("get incremented counter", func(t *testing.T) {
		for k, v := range list {
			val, err := st.GetCounter(ctx, k)
			assert.NoError(t, err)
			assert.Equal(t, *v.Delta, *val.Delta)
		}
	})
}

func TestMemStorage_UpdateBatch(t *testing.T) {
	tests := []struct {
		name        string
		metrics     []models.Metric
		expectError bool
	}{
		{
			name: "success",
			metrics: []models.Metric{
				{
					ID:    "test1",
					MType: models.TypeGauge,
					Value: new(float64),
				},
				{
					ID:    "test2",
					MType: models.TypeCounter,
					Delta: new(int64),
				},
			},
			expectError: false,
		},
		{
			name: "error",
			metrics: []models.Metric{
				{
					ID:    "test1",
					MType: "unknown",
					Value: new(float64),
				},
			},
			expectError: true,
		},
	}

	ctx := context.TODO()
	st := NewMemStorage("", false)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := st.UpdateBatch(ctx, test.metrics)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
