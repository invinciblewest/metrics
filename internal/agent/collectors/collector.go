package collectors

import (
	"context"
	"github.com/invinciblewest/metrics/pkg/worker"
)

type Collector interface {
	Collect(ctx context.Context) error
}

func CollectMetrics(workersPool *worker.Pool, collectors ...Collector) {
	for _, c := range collectors {
		workersPool.AddJob(func(ctx context.Context) error {
			return c.Collect(ctx)
		})
	}
}
