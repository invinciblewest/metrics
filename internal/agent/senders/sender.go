package senders

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/storage"
	"strconv"
)

type Sender interface {
	Send(mType string, mName string, mValue string) error
}

func SendMetrics(st *storage.MemStorage, senders ...Sender) error {
	fmt.Println("sending metrics to server...")
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
	fmt.Println("metrics have been sent")
	return nil
}
