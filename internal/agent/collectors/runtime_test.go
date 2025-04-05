package collectors

import (
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRuntimeCollector(t *testing.T) {
	st := memstorage.NewMemStorage("", false)
	c := NewRuntimeCollector(st)
	assert.Implements(t, (*Collector)(nil), c)
}

func TestRuntimeCollector_Collect(t *testing.T) {
	st := memstorage.NewMemStorage("", false)
	c := NewRuntimeCollector(st)

	t.Run("collect error", func(t *testing.T) {
		err := c.Collect()
		require.NoError(t, err)
	})

	t.Run("check gauges", func(t *testing.T) {
		gaugeKeyList := []string{
			"Alloc",
			"BuckHashSys",
			"Frees",
			"GCCPUFraction",
			"GCSys",
			"HeapAlloc",
			"HeapIdle",
			"HeapInuse",
			"HeapObjects",
			"HeapReleased",
			"HeapSys",
			"LastGC",
			"Lookups",
			"MCacheInuse",
			"MCacheSys",
			"MSpanInuse",
			"MSpanSys",
			"Mallocs",
			"NextGC",
			"NumForcedGC",
			"NumGC",
			"OtherSys",
			"PauseTotalNs",
			"StackInuse",
			"StackSys",
			"Sys",
			"TotalAlloc",
			"RandomValue",
		}

		gaugeList := st.GetGaugeList()
		for _, v := range gaugeKeyList {
			assert.Contains(t, gaugeList, v)
		}
	})
	t.Run("check counters", func(t *testing.T) {
		counterKeyList := []string{
			"PollCount",
		}

		counterList := st.GetCounterList()
		for _, v := range counterKeyList {
			assert.Contains(t, counterList, v)
		}
	})
}
