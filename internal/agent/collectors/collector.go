package collectors

import (
	"context"

	"github.com/invinciblewest/metrics/pkg/worker"
)

// Collector интерфейс для сбора метрик.
type Collector interface {
	// Collect собирает метрики и сохраняет их в хранилище.
	Collect(ctx context.Context) error
}

// CollectMetrics добавляет задачи сбора метрик в пул работников.
func CollectMetrics(workersPool *worker.Pool, collectors ...Collector) {
	for _, c := range collectors {
		workersPool.AddJob(func(ctx context.Context) error {
			return c.Collect(ctx)
		})
	}
}
