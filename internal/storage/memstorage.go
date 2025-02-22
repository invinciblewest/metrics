package storage

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (st *MemStorage) UpdateGauge(name string, value float64) {
	st.gauges[name] = value
}

func (st *MemStorage) UpdateCounter(name string, value int64) {
	st.counters[name] += value
}
