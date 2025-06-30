package services

import (
	"context"
	"errors"

	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
)

// MetricsService предоставляет методы для работы с метриками в хранилище.
type MetricsService struct {
	st storage.Storage
}

// NewMetricsService создает новый экземпляр MetricsService с заданным хранилищем.
func NewMetricsService(st storage.Storage) MetricsService {
	return MetricsService{
		st: st,
	}
}

// Update обновляет метрику в хранилище в зависимости от ее типа.
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

// UpdateBatch обновляет пакет метрик в хранилище.
func (ms *MetricsService) UpdateBatch(ctx context.Context, metrics []models.Metric) error {
	return ms.st.UpdateBatch(ctx, metrics)
}

// Get извлекает метрику из хранилища по типу и идентификатору.
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

// PingStorage проверяет доступность хранилища метрик.
func (ms *MetricsService) PingStorage(ctx context.Context) bool {
	if err := ms.st.Ping(ctx); err != nil {
		return false
	}
	return true
}
