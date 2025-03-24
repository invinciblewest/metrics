package storage

import (
	"errors"
	"github.com/invinciblewest/metrics/internal/models"
)

type MemStorage struct {
	gauges   GaugeList
	counters CounterList
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(GaugeList),
		counters: make(CounterList),
	}
}

func (st *MemStorage) UpdateGauge(metrics models.Metrics) error {
	if metrics.MType != models.TypeGauge {
		return errors.New("wrong type")
	}

	st.gauges[metrics.ID] = metrics
	return nil
}

func (st *MemStorage) GetGauge(id string) (models.Metrics, error) {
	value, exists := st.gauges[id]
	if exists {
		return value, nil
	} else {
		return models.Metrics{}, errors.New("not found")
	}
}

func (st *MemStorage) GetGaugeList() GaugeList {
	return st.gauges
}

func (st *MemStorage) UpdateCounter(metrics models.Metrics) error {
	if metrics.MType != models.TypeCounter {
		return errors.New("wrong type")
	}

	currentMetrics, exists := st.counters[metrics.ID]
	if exists {
		*metrics.Delta += *currentMetrics.Delta
	}

	st.counters[metrics.ID] = metrics
	return nil
}

func (st *MemStorage) GetCounter(id string) (models.Metrics, error) {
	value, exists := st.counters[id]
	if exists {
		return value, nil
	} else {
		return models.Metrics{}, errors.New("not found")
	}
}

func (st *MemStorage) GetCounterList() CounterList {
	return st.counters
}
