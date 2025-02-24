package storage

type GaugeList map[string]float64
type CounterList map[string]int64

type Storage interface {
	UpdateGauge(name string, value float64)
	GetGauge(name string) (float64, error)
	GetGaugeList() GaugeList
	UpdateCounter(name string, value int64)
	GetCounter(name string) (int64, error)
	GetCounterList() CounterList
}
