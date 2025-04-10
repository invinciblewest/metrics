package collectors

import "context"

type Collector interface {
	Collect(ctx context.Context) error
}

func CollectMetrics(ctx context.Context, collectors ...Collector) error {
	for _, c := range collectors {
		if err := c.Collect(ctx); err != nil {
			return err
		}
	}
	return nil
}
