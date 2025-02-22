package main

import (
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
	serverAddr     = "http://localhost:8080"
)

func main() {
	st := storage.NewMemStorage()
	runtimeCollector := collectors.NewRuntimeCollector()
	httpSender := senders.NewHttpSender(serverAddr, http.DefaultClient)

	for {
		if err := collectors.CollectMetrics(st, runtimeCollector); err != nil {
			panic(err)
		}

		if ((st.GetCounter("PollCount") * pollInterval) % reportInterval) == 0 {
			if err := senders.SendMetrics(st, httpSender); err != nil {
				panic(err)
			}
		}
		time.Sleep(pollInterval * time.Second)
	}
}
