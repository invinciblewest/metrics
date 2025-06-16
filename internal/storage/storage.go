package storage

import (
	"context"
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
	UpdateGauge(ctx context.Context, metric models.Metric) error
	GetGauge(ctx context.Context, id string) (models.Metric, error)
	GetGaugeList(ctx context.Context) GaugeList
	UpdateCounter(ctx context.Context, metric models.Metric) error
	GetCounter(ctx context.Context, id string) (models.Metric, error)
	GetCounterList(ctx context.Context) CounterList
	UpdateBatch(ctx context.Context, metrics []models.Metric) error
	Save(ctx context.Context) error
	Load(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}
