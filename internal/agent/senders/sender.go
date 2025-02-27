package senders

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"log"
	"strconv"
)

type Sender interface {
	Send(mType string, mName string, mValue string) error
}

func SendMetrics(st *storage.MemStorage, senders ...Sender) error {
	log.Println("sending metrics to server...")
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
	log.Println("metrics have been sent")
	return nil
}
