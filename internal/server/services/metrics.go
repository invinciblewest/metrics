package services

import (
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

func (ms *MetricsService) Update(metrics models.Metric) (models.Metric, error) {
	switch metrics.MType {
	case models.TypeGauge:
		if metrics.Value == nil {
			return metrics, errors.New("value is nil")
		}
		if err := ms.st.UpdateGauge(metrics); err != nil {
			return metrics, err
		}
	case models.TypeCounter:
		if metrics.Delta == nil {
			return metrics, errors.New("delta is nil")
		}
		if err := ms.st.UpdateCounter(metrics); err != nil {
			return metrics, err
		}
	default:
		return models.Metric{}, storage.ErrWrongType
	}

	return metrics, nil
}

func (ms *MetricsService) Get(mType, id string) (models.Metric, error) {
	var result models.Metric

	if id == "" {
		return result, errors.New("id is empty")
	}

	switch mType {
	case models.TypeGauge:
		value, err := ms.st.GetGauge(id)
		if err != nil {
			return result, storage.ErrNotFound
		}
		result = value
	case models.TypeCounter:
		value, err := ms.st.GetCounter(id)
		if err != nil {
			return result, err
		}
		result = value
	default:
		return result, storage.ErrWrongType
	}

	return result, nil
}
