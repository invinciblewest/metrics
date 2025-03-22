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

func (ms *MetricsService) Update(metrics models.Metrics) (models.Metrics, error) {
	if !metrics.CheckType() {
		return models.Metrics{}, errors.New("wrong type")
	}

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
	}

	return metrics, nil
}

func (ms *MetricsService) Get(mType, id string) (models.Metrics, error) {
	var result models.Metrics

	if id == "" || !models.CheckType(mType) {
		return result, errors.New("wrong type")
	}

	switch mType {
	case models.TypeGauge:
		value, err := ms.st.GetGauge(id)
		if err != nil {
			return result, errors.New("not found")
		}
		result = value
	case models.TypeCounter:
		value, err := ms.st.GetCounter(id)
		if err != nil {
			return result, errors.New("not found")
		}
		result = value
	}

	return result, nil
}
