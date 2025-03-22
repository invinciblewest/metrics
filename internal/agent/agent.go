package agent

import (
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"time"
)

type Agent struct {
	st         storage.Storage
	collectors []collectors.Collector
	senders    []senders.Sender
	pInterval  int
	rInterval  int
}

func NewAgent(
	st storage.Storage,
	collectors []collectors.Collector,
	senders []senders.Sender,
	pInterval int,
	rInterval int,
) *Agent {
	return &Agent{
		st:         st,
		collectors: collectors,
		senders:    senders,
		pInterval:  pInterval,
		rInterval:  rInterval,
	}
}

func (a *Agent) Run() error {
	for {
		if err := collectors.CollectMetrics(a.collectors...); err != nil {
			return err
		}

		pc, err := a.st.GetCounter("PollCount")
		if err != nil {
			return err
		}
		if ((int(*pc.Delta) * a.pInterval) % a.rInterval) == 0 {
			if err := senders.SendMetrics(a.st, a.senders...); err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(a.pInterval) * time.Second)
	}
}
