package storage

import "errors"

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

func (st *MemStorage) UpdateGauge(name string, value float64) {
	st.gauges[name] = value
}

func (st *MemStorage) GetGauge(name string) (float64, error) {
	value, exists := st.gauges[name]
	if exists {
		return value, nil
	} else {
		return 0, errors.New("not found")
	}
}

func (st *MemStorage) GetGaugeList() GaugeList {
	return st.gauges
}

func (st *MemStorage) UpdateCounter(name string, value int64) {
	st.counters[name] += value
}

func (st *MemStorage) GetCounterList() CounterList {
	return st.counters
}

func (st *MemStorage) GetCounter(name string) (int64, error) {
	value, exists := st.counters[name]
	if exists {
		return value, nil
	} else {
		return 0, errors.New("not found")
	}
}
