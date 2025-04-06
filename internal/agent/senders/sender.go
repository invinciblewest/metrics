package senders

import (
	"context"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
)

type Sender interface {
	SendMetric(ctx context.Context, metrics []models.Metric) error
}

func SendMetrics(ctx context.Context, st storage.Storage, senders ...Sender) error {
	logger.Log.Info("sending metrics to server...")
	for _, s := range senders {
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
	}
	logger.Log.Info("metrics have been sent")
	return nil
}
