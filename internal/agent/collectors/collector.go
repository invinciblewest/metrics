package collectors

import (
	"github.com/invinciblewest/metrics/internal/storage"
)

type Collector interface {
	Collect(st *storage.MemStorage) error
}

func CollectMetrics(st *storage.MemStorage, collectors ...Collector) error {
	for _, c := range collectors {
		if err := c.Collect(st); err != nil {
			return err
		}
	}
	return nil
}
