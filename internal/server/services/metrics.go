package services

import (
	"context"
	"errors"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
)

type MetricsService struct {
	st storage.Storage
}

func NewMetricsService(st storage.Storage) MetricsService {
	return MetricsService{
		st: st,
	}
}

func (ms *MetricsService) Update(ctx context.Context, metrics models.Metric) (models.Metric, error) {
	switch metrics.MType {
	case models.TypeGauge:
		if metrics.Value == nil {
			return metrics, errors.New("value is nil")
		}
		if err := ms.st.UpdateGauge(ctx, metrics); err != nil {
			return metrics, err
		}
	case models.TypeCounter:
		if metrics.Delta == nil {
			return metrics, errors.New("delta is nil")
		}
		if err := ms.st.UpdateCounter(ctx, metrics); err != nil {
			return metrics, err
		}
	default:
		return models.Metric{}, storage.ErrWrongType
	}

	return metrics, nil
}

func (ms *MetricsService) UpdateBatch(ctx context.Context, metrics []models.Metric) error {
	return ms.st.UpdateBatch(ctx, metrics)
}

func (ms *MetricsService) Get(ctx context.Context, mType, id string) (models.Metric, error) {
	var result models.Metric

	if id == "" {
		return result, errors.New("id is empty")
	}

	switch mType {
	case models.TypeGauge:
		value, err := ms.st.GetGauge(ctx, id)
		if err != nil {
			return result, storage.ErrNotFound
		}
		result = value
	case models.TypeCounter:
		value, err := ms.st.GetCounter(ctx, id)
		if err != nil {
			return result, storage.ErrNotFound
		}
		result = value
	default:
		return result, storage.ErrWrongType
	}

	return result, nil
}

func (ms *MetricsService) PingStorage(ctx context.Context) bool {
	if err := ms.st.Ping(ctx); err != nil {
		return false
	}
	return true
}
