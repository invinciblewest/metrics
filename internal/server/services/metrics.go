package services

import (
	"errors"
	"github.com/invinciblewest/metrics/internal/storage"
	"strconv"
)

type MetricsService struct {
	st storage.Storage
}

func NewMetricsService(st storage.Storage) MetricsService {
	return MetricsService{
		st: st,
	}
}

func (ms *MetricsService) checkType(metricType string) bool {
	return metricType == "gauge" || metricType == "counter"
}

func (ms *MetricsService) Update(metricType, metricName, metricValue string) error {
	if !ms.checkType(metricType) {
		return errors.New("wrong type")
	}

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		ms.st.UpdateGauge(metricName, value)
	default:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		ms.st.UpdateCounter(metricName, value)
	}

	return nil
}

func (ms *MetricsService) GetString(metricType, metricName string) (string, error) {
	if metricName == "" || !ms.checkType(metricType) {
		return "", errors.New("wrong type")
	}

	var result string

	switch metricType {
	case "gauge":
		value, err := ms.st.GetGauge(metricName)
		if err != nil {
			return result, errors.New("not found")
		}
		result = strconv.FormatFloat(value, 'f', -1, 64)
	default:
		value, err := ms.st.GetCounter(metricName)
		if err != nil {
			return result, errors.New("not found")
		}
		result = strconv.FormatInt(value, 10)
	}

	return result, nil
}
