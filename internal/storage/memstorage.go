package storage

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

func (st *MemStorage) GetGauge(name string) float64 {
	return st.gauges[name]
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

func (st *MemStorage) GetCounter(name string) int64 {
	return st.counters[name]
}
