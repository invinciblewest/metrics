package senders

import (
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
)

type Sender interface {
	Send(metrics models.Metrics) error
}

func SendMetrics(st storage.Storage, senders ...Sender) error {
	logger.Log.Info("sending metrics to server...")
	for _, s := range senders {
		for _, v := range st.GetGaugeList() {
			if err := s.Send(v); err != nil {
				return err
			}
		}
		for _, v := range st.GetCounterList() {
			if err := s.Send(v); err != nil {
				return err
			}
		}
	}
	logger.Log.Info("metrics have been sent")
	return nil
}
