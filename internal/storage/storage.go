package storage

import (
	"errors"
	"github.com/invinciblewest/metrics/internal/models"
)

type GaugeList map[string]models.Metric
type CounterList map[string]models.Metric

var (
	ErrNotFound  = errors.New("not found")
	ErrWrongType = errors.New("wrong type")
)

type Storage interface {
	UpdateGauge(metric models.Metric) error
	GetGauge(id string) (models.Metric, error)
	GetGaugeList() GaugeList
	UpdateCounter(metric models.Metric) error
	GetCounter(id string) (models.Metric, error)
	GetCounterList() CounterList
	Save() error
	Load() error
	Ping() error
	Close()
}
