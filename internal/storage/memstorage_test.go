package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	st := NewMemStorage()
	assert.Implements(t, (*Storage)(nil), st)
}

func TestMemStorage_Gauge(t *testing.T) {
	st := NewMemStorage()

	list := GaugeList{
		"test1": 3.14,
		"test2": 14.3,
	}
	t.Run("update gauge", func(t *testing.T) {
		for k, v := range list {
			st.UpdateGauge(k, v)
		}
	})
	t.Run("get gauge", func(t *testing.T) {
		for k, v := range list {
			assert.Equal(t, v, st.GetGauge(k))
		}
	})
	t.Run("get gauge list", func(t *testing.T) {
		assert.Equal(t, list, st.GetGaugeList())
	})
}

func TestMemStorage_Counter(t *testing.T) {
	st := NewMemStorage()

	list := CounterList{
		"test1": 1,
		"test2": 2,
	}
	t.Run("update counter", func(t *testing.T) {
		for k, v := range list {
			st.UpdateCounter(k, v)
		}
	})
	t.Run("get counter", func(t *testing.T) {
		for k, v := range list {
			assert.Equal(t, v, st.GetCounter(k))
		}
	})
	t.Run("get counter list", func(t *testing.T) {
		assert.Equal(t, list, st.GetCounterList())
	})
	t.Run("increment counter", func(t *testing.T) {
		for k, v := range list {
			st.UpdateCounter(k, v)
		}
	})
	t.Run("get counter after increment", func(t *testing.T) {
		for k, v := range list {
			assert.Equal(t, v*2, st.GetCounter(k))
		}
	})
}
