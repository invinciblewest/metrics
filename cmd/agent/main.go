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
	collectorsList := []collectors.Collector{
		collectors.NewRuntimeCollector(),
	}
	sendersList := []senders.Sender{
		senders.NewHTTPSender(serverAddr, http.DefaultClient),
	}

	runAgent(st, collectorsList, sendersList, pollInterval, reportInterval)
}

func runAgent(
	st *storage.MemStorage,
	c []collectors.Collector,
	s []senders.Sender,
	pInterval int,
	rInterval int,
) {
	for {
		if err := collectors.CollectMetrics(st, c...); err != nil {
			panic(err)
		}

		pc, _ := st.GetCounter("PollCount")
		if ((int(pc) * pInterval) % rInterval) == 0 {
			if err := senders.SendMetrics(st, s...); err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Duration(pInterval) * time.Second)
	}
}
