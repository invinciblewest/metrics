package collectors

type Collector interface {
	Collect() error
}

func CollectMetrics(collectors ...Collector) error {
	for _, c := range collectors {
		if err := c.Collect(); err != nil {
			return err
		}
	}
	return nil
}
