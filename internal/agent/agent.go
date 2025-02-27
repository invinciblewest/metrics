package agent

import (
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"time"
)

type Agent struct {
	st        *storage.MemStorage
	c         []collectors.Collector
	s         []senders.Sender
	pInterval int
	rInterval int
}

func NewAgent(
	st *storage.MemStorage,
	c []collectors.Collector,
	s []senders.Sender,
	pInterval int,
	rInterval int,
) *Agent {
	return &Agent{
		st:        st,
		c:         c,
		s:         s,
		pInterval: pInterval,
		rInterval: rInterval,
	}
}

func (a *Agent) Run() error {
	for {
		if err := collectors.CollectMetrics(a.st, a.c...); err != nil {
			return err
		}

		pc, _ := a.st.GetCounter("PollCount")
		if ((int(pc) * a.pInterval) % a.rInterval) == 0 {
			if err := senders.SendMetrics(a.st, a.s...); err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(a.pInterval) * time.Second)
	}
}
