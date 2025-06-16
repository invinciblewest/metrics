package senders

import (
	"context"

	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/invinciblewest/metrics/pkg/worker"
	"go.uber.org/zap"
)

// Sender интерфейс для отправки метрик на сервер.
type Sender interface {
	// SendMetric отправляет список метрик на сервер.
	SendMetric(ctx context.Context, metrics []models.Metric) error
}

// SendMetrics отправляет метрики на сервер с использованием пула воркеров и заданных отправителей.
func SendMetrics(workersPool *worker.Pool, st storage.Storage, senders ...Sender) {
	for _, s := range senders {
		workersPool.AddJob(func(ctx context.Context) error {
			logger.Log.Info("sending metrics to server...")

			metrics := make([]models.Metric, 0)
			for _, v := range st.GetGaugeList(ctx) {
				metrics = append(metrics, v)
			}
			for _, v := range st.GetCounterList(ctx) {
				metrics = append(metrics, v)
			}
			err := s.SendMetric(ctx, metrics)
			if err != nil {
				logger.Log.Error("failed to send metrics: ", zap.Error(err))
				return err
			}

			logger.Log.Info("metrics have been sent")
			return nil
		})
	}
}
