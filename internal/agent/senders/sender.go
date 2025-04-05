package senders

import (
	"context"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
)

type Sender interface {
	SendMetric(ctx context.Context, metric models.Metric) error
}

func SendMetrics(ctx context.Context, st storage.Storage, senders ...Sender) error {
	logger.Log.Info("sending metrics to server...")
	for _, s := range senders {
		for _, v := range st.GetGaugeList(ctx) {
			if err := s.SendMetric(ctx, v); err != nil {
				return err
			}
		}
		for _, v := range st.GetCounterList(ctx) {
			if err := s.SendMetric(ctx, v); err != nil {
				return err
			}
		}
	}
	logger.Log.Info("metrics have been sent")
	return nil
}
