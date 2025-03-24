package storage

import "github.com/invinciblewest/metrics/internal/models"

type GaugeList map[string]models.Metrics
type CounterList map[string]models.Metrics

type Storage interface {
	UpdateGauge(metrics models.Metrics) error
	GetGauge(id string) (models.Metrics, error)
	GetGaugeList() GaugeList
	UpdateCounter(metrics models.Metrics) error
	GetCounter(id string) (models.Metrics, error)
	GetCounterList() CounterList
}
