package agent

import (
	"context"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/invinciblewest/metrics/pkg/worker"
	"go.uber.org/zap"
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

func (a *Agent) Run(ctx context.Context, rateLimit int) error {
	pollTicker := time.NewTicker(time.Duration(a.pInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(a.rInterval) * time.Second)
	defer reportTicker.Stop()

	workersPool := worker.NewPool()
	defer workersPool.Stop()

	errorsCh := workersPool.Start(ctx, rateLimit)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("context done")
				return
			case <-pollTicker.C:
				collectors.CollectMetrics(workersPool, a.collectors...)
			case <-reportTicker.C:
				senders.SendMetrics(workersPool, a.st, a.senders...)
			}
		}
	}()

	for err := range errorsCh {
		if err != nil {
			logger.Log.Error("worker error", zap.Error(err))
			return err
		}
	}

	return nil
}
