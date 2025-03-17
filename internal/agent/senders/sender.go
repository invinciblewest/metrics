package senders

import (
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/storage"
	"strconv"
)

type Sender interface {
	Send(mType string, mName string, mValue string) error
}

func SendMetrics(st storage.Storage, senders ...Sender) error {
	logger.Log.Info("Sending metrics to server...")
	for _, s := range senders {
		for k, v := range st.GetGaugeList() {
			val := strconv.FormatFloat(v, 'f', -1, 64)
			if err := s.Send("gauge", k, val); err != nil {
				return err
			}
		}
		for k, v := range st.GetCounterList() {
			val := strconv.FormatInt(v, 10)
			if err := s.Send("counter", k, val); err != nil {
				return err
			}
		}
	}
	logger.Log.Info("metrics have been sent")
	return nil
}
